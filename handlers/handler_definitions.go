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
	"tacos",
	"enchiladas",
	"burritos",
	"churros",
	"beans",
	"salsa",
	"nachos",
}

// TacoHandler responds to a message containing the substring "taco"
var TacoHandler = &HandlerDefinition{
	Match: func(message string) bool {
		return strings.Contains(strings.Join(goodFoods, " "), strings.ToLower(message))
	},
	Handle: func(reply ReplyFn, event *slack.MessageEvent) {
		goodFood, found := funk.FindString(goodFoods, func(food string) bool {
			if strings.Contains(food, event.Text) {
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

		subject := expression.FindAllString(event.Text, -1)

		if len(subject) == 0 {
			reply(event, "Sorry, I'm a little confused about what I am tracking karma on!")
		} else {
			buf := fmt.Sprintf("Tracking karma on %s", strings.Join(subject, ", "))
			reply(event, buf)
		}
	},
}
