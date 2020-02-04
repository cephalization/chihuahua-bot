package main

import (
	"log"
	"os"
	"strings"

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

	auth, err := api.AuthTest()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	botUser, err := api.GetUserInfo(auth.UserID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	botID := botUser.ID

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Printf("\n\nChihuahua bot is running, %s\n", ansi.Color("bark! bark!", "green"))

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// Ignore messages from this bot
			if ev.BotID != "" || ev.BotID == botID {
				break
			}

			if strings.Contains(strings.ToLower(ev.Text), "taco") {
				rtm.PostMessage(ev.Channel, slack.MsgOptionText("mmm... tacos...", false))
			}
		default:
		}
	}
}
