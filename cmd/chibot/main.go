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

	// Setup an error context for picking up errors with the bot thread
	botErrorContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Instantiate slack bot
	bot := slacker.NewClient(token)

	// Register commands onto the bot
	commands.RegisterCommands(bot)

	// Configure bot use the error context when it encounters errors in the background
	err = bot.Listen(botErrorContext)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("Chihuahua bot is running %s!", ansi.Color("bark", ansi.Cyan))
}
