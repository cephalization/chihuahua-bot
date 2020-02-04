package main

import (
	"context"
	"log"
	"os"

	"github.com/cephalization/chihuahua-bot/commands"
	"github.com/cephalization/chihuahua-bot/utils"

	"github.com/mgutz/ansi"
	"github.com/shomali11/slacker"
)

func main() {
	// Grab API token from env variables
	token, err := utils.GetEnv("token")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Instantiate slack bot
	bot := slacker.NewClient(token, slacker.WithDebug(true))

	// Register commands onto the bot
	commands.RegisterCommands(bot)

	log.Printf("\n\nChihuahua bot is running, %s\n", ansi.Color("bark! bark!", "green"))

	// Setup a context for picking up messages from the bot thread
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configure bot to listen for messages, bot is now running
	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
