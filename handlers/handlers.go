package handlers

import "github.com/nlopes/slack"

// ReplyFn replies to the channel that triggered this event with a message
type ReplyFn func(*slack.MessageEvent, string)

// Handler contains utils to help message handlers respond to messages
type Handler struct {
	Reply ReplyFn
	RTM   *slack.RTM
}

var definitions = [...]*HandlerDefinition{
	TacoHandler,
	KarmaHandler,
}

// HandleMessages from the real time messaging api, passing them off to the correct fn
func (handler *Handler) HandleMessages(event *slack.MessageEvent) {
	message := event.Text

	for _, definition := range definitions {
		if definition.Match(message) {
			handler.RTM.NewTypingMessage(event.Channel)

			// This function must reply with something or else the bot will appear to be typing forever
			definition.Handle(handler.Reply, event)
		}
	}
}
