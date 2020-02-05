package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/thoas/go-funk"
)

// HandlerDefinition models a function that replies to a message that matches its criteria
type HandlerDefinition struct {
	Match  func(string) bool
	Handle func(ReplyFn, *slack.MessageEvent)
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
	Handle: func(reply ReplyFn, event *slack.MessageEvent) {
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
	Handle: func(reply ReplyFn, event *slack.MessageEvent) {
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
	Handle: func(reply ReplyFn, event *slack.MessageEvent) {
		subjects := strings.Split(event.Text, " ")[1:]

		buf := ""
		for _, subject := range subjects {
			karma := "?"
			buf += fmt.Sprintf("Karma for `%s` is %s\n", subject, karma)
		}

		reply(event, buf)
	},
}