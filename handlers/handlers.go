package handlers

import "github.com/nlopes/slack"

// Handler contains utils to help message handlers respond to messages
type Handler struct {
	Reply func(event *slack.MessageEvent, message string)
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
			definition.Handle(handler.Reply, event)
		}
	}
}
