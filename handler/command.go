package handler

import traqbot "github.com/motoki317/traq-bot"

type command struct {
	name   string
	help   string
	handle func(h *handler, payload *traqbot.MessageCreatedPayload, args []string) error
}
