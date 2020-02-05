package main

import (
	"log"
	"os"

	"github.com/cephalization/chihuahua-bot/handlers"
	"github.com/cephalization/chihuahua-bot/utils"

	"github.com/mgutz/ansi"
)

func main() {
	// Grab API token from env variables
	token, err := utils.GetEnv("token")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	server, err := handlers.NewClient(token)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("\n\nChihuahua bot is running, %s\n", ansi.Color("bark! bark!", "green"))

	server.Listen()
}
