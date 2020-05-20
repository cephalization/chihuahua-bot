package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cephalization/chihuahua-bot/utils"
	"github.com/cephalization/minestat/Go/minestat"
	"github.com/nlopes/slack"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HandlerDefinition models a function that replies to a message that matches its criteria
type HandlerDefinition struct {
	Match  func(string) bool
	Handle func(ReplyFn, *slack.MessageEvent, *Handler)
}

var goodFoods = []string{
	"taco",
	"enchilada",
	"burrito",
	"churro",
	"bean",
	"salsa",
	"nacho",
}

// The assorted noises the ting can make
var noisesTheTingMakes = []string{
	"skrrrrrrrrrr",
	"braaaaaat",
	"pop pop",
	"kun kun",
	"doon doon",
	"skyaaaaaaaaaaaaaaaaaaaaaa",
	"BOOM",
	"BRAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT",
	"skrip skirp skirp",
	"draaaaaa",
}

var adjectives []string

// ThingGotRegex - for use in ThingGotHandler()
var ThingGotRegex = "(\\w+) (?i)got (?i)(a |that |an |dat |them |some |sum |its )?(.*)\\?"

// ThingGotHandler responds to a message that contain a substring that matches
// the regex "(\w+) got (\w+)?"
var ThingGotHandler = &HandlerDefinition{
	Match: func(message string) bool {
		matched, _ := regexp.MatchString(ThingGotRegex, message)
		return matched
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		matchRE := regexp.MustCompile(ThingGotRegex)
		match := matchRE.FindStringSubmatch(event.Text)

		if len(adjectives) == 0 {
			// Populate adjectives list from disk on file
			adjectives = utils.PopulateAdjectivesFromFile()
		}

		// 'match' should contain 3 elements: the first being the entire
		// matched string, the next two being the two capture groups.
		if len(match) < 3 {
			return
		}

		// Index to a random spot in the list of adjectives
		rand.Seed(time.Now().Unix())
		index := rand.Intn(len(adjectives) - 1)
		adjective := adjectives[index]
		aestheticAdjective := strings.Join(strings.Split(adjective, ""), " ")

		// Handle prepositions
		var substr string
		if len(match) == 4 {
			substr = fmt.Sprintf("%s*%s* %s", match[2],
				aestheticAdjective, match[3])
		} else {
			substr = fmt.Sprintf("*%s* %s", adjective,
				match[2])
		}

		// Replace 'i'or 'I' with the user's name
		if strings.ToLower(match[1]) == "i" {
			user, _ := handler.API.GetUserInfo(event.User)
			match[1] = user.Profile.RealName

			// Replace instances of 'my' with 'their'
			myRE := regexp.MustCompile("(?i)(my)")
			substr = myRE.ReplaceAllLiteralString(substr, "their")
		}

		// Replace "chibot" with I
		if strings.ToLower(match[1]) == "chibot" {
			match[1] = "I"
		}

		buf := fmt.Sprintf("%s got %s\n", match[1], substr)
		reply(event, buf)
	},
}

// BigShaqHandler responds to messages containing the string "ting go?"
var BigShaqHandler = &HandlerDefinition{
	Match: func(message string) bool {
		matched, _ := regexp.MatchString("ting go\\?", message)
		return matched
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		// Determine a number of noises the ting makes...the number
		// 15 is kind of arbitrary, anything longer would be even _more_
		// obnoxious.
		rand.Seed(time.Now().Unix())
		length := rand.Intn(15)

		buf := "ting go"

		// Append noises to 'buf'
		for i := 0; i < length; i++ {
			index := rand.Intn(len(noisesTheTingMakes) - 2)
			buf += fmt.Sprintf(" %s", noisesTheTingMakes[index])
		}
		reply(event, buf)
	},
}

// TacoHandler responds to a message containing the substring "taco"
var TacoHandler = &HandlerDefinition{
	Match: func(message string) bool {
		validFood := false

		for _, food := range goodFoods {
			if strings.Contains(strings.ToLower(message), food) {
				validFood = true
				continue
			}
		}

		return validFood
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		goodFood, found := funk.FindString(goodFoods, func(food string) bool {
			if strings.Contains(event.Text, food) {
				return true
			}

			return false
		})

		if !found {
			reply(event, "mmm... food...")
			return
		}

		buf := fmt.Sprintf("mmm... %s...", goodFood)
		reply(event, buf)
	},
}

// Struct to store query information for a mongo query
// bson.M is just a map literal, bson.M{"key": value}
type update struct {
	filter bson.M // map of filter query info
	update bson.M // map of update operations
}

var karmaRegex = `\w*\b(--|\+\+)`

// KarmaHandler responds to a message ending in ++ or -- by tracking score on the word in the message
var KarmaHandler = &HandlerDefinition{
	Match: func(message string) bool {
		if len(message) <= 2 {
			return false
		}

		matched, err := regexp.Match(karmaRegex, []byte(message))
		if err != nil {
			return false
		}

		return matched
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		expression, err := regexp.Compile(karmaRegex)
		if err != nil {
			reply(event, "Could not track karma on your message, sorry!")
			return
		}

		// array of all strings that match the regex
		// ex. 'apple++'
		subjects := expression.FindAllString(event.Text, -1)
		for i, s := range subjects {
			subjects[i] = strings.ToLower(s)
		}

		const add = byte('+')
		const subtract = byte('-')

		toUpdate := make(map[string]int, 0)
		var updates []*update

		// If we can get the user's info, try and check they don't upvote themselves later
		ignoreUserCheck := false
		userName := ""
		user, err := handler.API.GetUserInfo(event.User)
		if err != nil {
			ignoreUserCheck = true
		} else {
			userName = user.Profile.RealName
		}
		chibotDelta := 0

		// Gather scores
		// Since a user can karma the same thing multiple times we use a map
		for _, subject := range subjects {
			action := subject[len(subject)-1]
			subject := subject[:len(subject)-2]
			skip := false

			// If the subject appears in a part of the user's name, ignore this subject
			if !ignoreUserCheck && action == add {
				for _, s := range strings.Split(userName, " ") {
					if strings.Contains(strings.ToLower(s), strings.ToLower(subject)) {
						skip = true
						break
					}
				}
			}

			if skip {
				continue
			}

			if action == add {
				toUpdate[subject]++
				if strings.ToLower(subject) == "chibot" {
					chibotDelta++
				}
			} else if action == subtract {
				toUpdate[subject]--
				if strings.ToLower(subject) == "chibot" {
					chibotDelta--
				}
			}
		}

		// Build update structs for mongo queries
		for name, score := range toUpdate {
			updates = append(
				updates,
				&update{
					filter: bson.M{"name": name},
					update: bson.M{"$inc": bson.M{"score": score}},
				})
		}

		// Make update queries
		emoji := "heavy_check_mark"
		if chibotDelta > 0 {
			emoji = "sdab"
		} else if chibotDelta < 0 {
			emoji = "cry"
		}

		opts := options.Update().SetUpsert(true) // we will create subject in db if not found
		collection := handler.DB.Collection("karma")
		for _, u := range updates {
			u := u

			_, err := collection.UpdateOne(context.TODO(), u.filter, u.update, opts)
			if err != nil {
				emoji = "bangbang"
				continue
			}
		}

		// reply(event, buf)
		handler.API.AddReaction(emoji, slack.NewRefToMessage(event.Channel, event.Timestamp))
	},
}

// ShowKarmaHandler reports the store karma value of everything after the word "karma"
var ShowKarmaHandler = &HandlerDefinition{
	Match: func(message string) bool {
		return strings.HasPrefix(strings.ToLower(message), "karma") && len(strings.TrimSpace(message)) > len("karma")
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		subjects := funk.UniqString(strings.Split(strings.ToLower(event.Text), " ")[1:])

		buf := ""
		collection := handler.DB.Collection("karma")
		for _, subject := range subjects {
			subject := subject

			var result bson.M
			err := collection.FindOne(context.TODO(), bson.M{"name": subject}).Decode(&result)
			if err != nil {
				buf += fmt.Sprintf("No karma found for `%s`\n", subject)
				continue
			}

			karma := result["score"]

			buf += fmt.Sprintf("Karma for `%s` is `%d`\n", subject, karma)
		}

		reply(event, buf)
	},
}

var MinecraftHandler = &HandlerDefinition{
	Match: func(message string) bool {
		return strings.Contains(strings.ToLower(message), "minecraft server up") || strings.Contains(strings.ToLower(message), "server up")
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		noConnection := func() {
			handler.API.AddReaction("no_entry_sign", slack.NewRefToMessage(event.Channel, event.Timestamp))
		}
		addressesString, err := utils.GetEnv("MINECRAFT_ADDRESSES")
		if err != nil {
			noConnection()
			return
		}

		addresses := strings.Split(addressesString, ",")

		buf := ""

		for i, address := range addresses {
			addressParts := strings.Split(address, ":")
			if len(addressParts) != 2 {
				continue
			}

			minestat.Init(addressParts[0], addressParts[1], 2)

			primaryServer := "Primary "
			server := ""
			if i != 0 {
				server = strconv.Itoa(i + 1)
				primaryServer = ""
			}

			if minestat.Online {
				buf += fmt.Sprintf(`
--------------------------
*%sserver %s*
Address: %s:%s
Players: %s/%s
Version: %s
%s
__________________________
				`,
					primaryServer,
					server,
					minestat.Address,
					minestat.Port,
					minestat.Current_players,
					minestat.Max_players,
					minestat.Version,
					minestat.Motd,
				)
			}
		}

		if buf == "" {
			noConnection()
			return
		}

		reply(event, buf)
	},
}
