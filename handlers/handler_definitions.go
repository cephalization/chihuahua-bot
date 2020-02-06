package handlers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/thoas/go-funk"
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

// ThingGotRegex - for use in ThingGotHandler()
var ThingGotRegex = "(\\w+) got (a |that |an |dat )?(.*)\\?"

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

		// 'match' should contain 3 elements: the first being the entire
		// matched string, the next two being the two capture groups.
		if len(match) < 3 {
			return
		}

		// Index to a random spot in the list of adjectives
		rand.Seed(time.Now().Unix())
		index := rand.Intn(len(adjectives) - 1)

		// Handle prepositions
		var substr string
		if len(match) == 4 {
			substr = fmt.Sprintf("%s%s %s", match[2],
				adjectives[index], match[3])
		} else {
			substr = fmt.Sprintf("%s %s", adjectives[index],
				match[2])
		}

		// Replace 'i'or 'I' with the user's name
		if match[1] == "i" || match[1] == "I" {
			user, _ := handler.API.GetUserInfo(event.User)
			match[1] = user.Profile.RealName

			// Replace instances of 'my' with 'their'
			myRE := regexp.MustCompile("(?i)(my)")
			substr = myRE.ReplaceAllLiteralString(substr, "their")
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

var karmaRegex = `\w*\b(--|\+\+)`

// KarmaHandler responds to a message ending in ++ or -- by tracking score on the word in the message
// TODO: track karma in db
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

		const add = byte('+')
		const subtract = byte('-')

		// right now buffer is just gonna store a string saying what things to update
		// in the future it should probably be some kind of map so we can batch all the adds and subtracts
		// into 1-2 calls to the database
		buf := ""

		for _, subject := range subjects {
			action := subject[len(subject)-1]

			if action == add {
				// debug message
				buf += fmt.Sprintf("Adding karma to %s\n", subject[:len(subject)-2])

				// TODO: connect to db, increment score for this subject
			} else if action == subtract {
				// debug message
				buf += fmt.Sprintf("Removing karma from %s\n", subject[:len(subject)-2])

				// TODO: connect to db, decrement score for this subject
			}
		}

		if len(buf) > 0 {
			reply(event, buf)
		} else {
			reply(event, "Sorry, I'm a little confused about what I am tracking karma on!")
		}
	},
}

// ShowKarmaHandler reports the store karma value of everything after the word "karma"
var ShowKarmaHandler = &HandlerDefinition{
	Match: func(message string) bool {
		return strings.HasPrefix(strings.ToLower(message), "karma") && len(strings.TrimSpace(message)) > len("karma")
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent, handler *Handler) {
		subjects := strings.Split(event.Text, " ")[1:]

		buf := ""
		for _, subject := range subjects {
			karma := "?"
			buf += fmt.Sprintf("Karma for `%s` is %s\n", subject, karma)
		}

		reply(event, buf)
	},
}
