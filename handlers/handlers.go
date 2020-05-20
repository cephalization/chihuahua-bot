package handlers

import (
	"github.com/nlopes/slack"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReplyFn replies to the channel that triggered this event with a message
type ReplyFn func(*slack.MessageEvent, string)

// Handler contains utils to help message handlers respond to messages
type Handler struct {
	Reply ReplyFn
	RTM   *slack.RTM
	API   *slack.Client
	ID    string
	DB    *mongo.Database
}

// These definitions will listen for message text and then perform an action
var definitions = []*HandlerDefinition{
	TacoHandler,
	KarmaHandler,
	ShowKarmaHandler,
	ThingGotHandler,
	BigShaqHandler,
	MinecraftHandler,
}

// HandleMessages from the real time messaging api, passing them off to the correct fn
func (handler *Handler) HandleMessages(event *slack.MessageEvent) {
	message := event.Text

	for _, definition := range definitions {
		if definition.Match(message) {
			handler.RTM.SendMessage(handler.RTM.NewTypingMessage(event.Channel))

			// This function must reply with something or else the bot will appear to be typing forever
			definition.Handle(handler.Reply, event, handler)
		}
	}
}

// Listen infinitely for slack events (messages) in real time
func (handler *Handler) Listen() {
	handleMessages := handler.HandleMessages
	rtm := handler.RTM
	ID := handler.ID

	go rtm.ManageConnection()

	// Read messages as they come in, pass them to the appropriate handlers
	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			// Ignore messages from this bot
			if event.BotID != "" || event.BotID == ID {
				break
			}

			handleMessages(event)
			break
		default:
			break
		}
	}
}

// NewClient returns a slack client that is ready to listen to real time messages
func NewClient(token string, DB *mongo.Database) (*Handler, error) {
	handler := &Handler{}

	// slack client
	api := slack.New(token)

	// check validity of auth token
	auth, err := api.AuthTest()
	if err != nil {
		return handler, err
	}

	// Parse bot's auth info
	botUser, err := api.GetUserInfo(auth.UserID)
	if err != nil {
		return handler, err
	}

	// Our bot's ID
	botID := botUser.ID

	// Connect to the real time messaging socket
	// This gets us an async channel of messages
	rtm := api.NewRTM()

	// Create a closure around our ideal rtm messaging function
	// Use this for easy replying
	reply := func(event *slack.MessageEvent, message string) {
		rtm.PostMessage(event.Channel, slack.MsgOptionText(message, false))
	}

	// Create our client handler
	handler = &Handler{
		API:   api,
		RTM:   rtm,
		ID:    botID,
		Reply: reply,
		DB:    DB,
	}

	return handler, nil
}
