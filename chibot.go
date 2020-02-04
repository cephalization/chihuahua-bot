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
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Printf("\n\nChihuahua bot is running, %s\n", ansi.Color("bark! bark!", "green"))

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			log.Printf("\n\nTRIGGERED: %s\n\n", ev.Text)
			if strings.Contains(ev.Text, "taco") {
				rtm.PostMessage(ev.Channel, slack.MsgOptionText("mmm... tacos...", false))
			}
		default:
		}
	}
}
