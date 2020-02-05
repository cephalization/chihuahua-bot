package main

import (
	"log"
	"os"

	"github.com/cephalization/chihuahua-bot/handlers"
	"github.com/cephalization/chihuahua-bot/utils"

	"github.com/mgutz/ansi"
	"github.com/nlopes/slack"
)

func main() {
	// Grab API token from env variables
	token, err := utils.GetEnv("token")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	api := slack.New(token, slack.OptionDebug(true))

	// check validity of auth token
	auth, err := api.AuthTest()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Parse bot's auth info
	botUser, err := api.GetUserInfo(auth.UserID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	botID := botUser.ID

	// Connect to the real time messaging socket
	// This gets us an async channel of messages
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Printf("\n\nChihuahua bot is running, %s\n", ansi.Color("bark! bark!", "green"))

	// Create a closure around our ideal rtm messaging function
	// Use this for easy replying
	reply := func(event *slack.MessageEvent, message string) {
		rtm.PostMessage(event.Channel, slack.MsgOptionText(message, false))
	}

	// Instantiate a message handler with our reply fn
	messageHandler := handlers.Handler{Reply: reply}

	// Read messages as they come in, pass them to the appropriate handlers
	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			// Ignore messages from this bot
			if event.BotID != "" || event.BotID == botID {
				break
			}

			messageHandler.HandleMessages(event)
			break
		default:
			break
		}
	}
}
