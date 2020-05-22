package handler

import (
	"fmt"
	traqbot "github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	"github.com/motoki317/traq-message-indexer/repository"
	"log"
	"strings"
)

const prefix = "/"

type handler struct {
	repo     repository.Repository
	commands []command
}

func MessageReceived(repo repository.Repository) func(payload *traqbot.MessageCreatedPayload) {
	commands := []command{
		commandIndex(),
		commandSearch(),
	}
	// append help command
	commands = append(commands, command{})
	commands[len(commands)-1] = commandHelp(commands)

	h := handler{
		repo:     repo,
		commands: commands,
	}

	return func(payload *traqbot.MessageCreatedPayload) {
		h.handleMessage(payload)
	}
}

func (h *handler) handleMessage(payload *traqbot.MessageCreatedPayload) {
	fields := strings.Fields(payload.Message.PlainText)
	for i, field := range fields {
		if strings.HasPrefix(field, prefix) {
			// Check for each command match
			for _, command := range h.commands {
				if strings.ToLower(field) == prefix+command.name {
					// Command match
					h.handleCommand(payload, command, fields[i:])
				}
			}
		}
	}
}

func (h *handler) handleCommand(payload *traqbot.MessageCreatedPayload, command command, args []string) {
	log.Printf("[%s]<%s>: %s\n", payload.Message.ChannelID, payload.Message.User.Name, payload.Message.PlainText)
	if err := command.handle(h, payload, args); err != nil {
		// error on command process
		if _, err := api.PostMessage(payload.Message.ChannelID, fmt.Sprintf("An error occurred while processing your command\n%s", err)); err != nil {
			log.Println(err)
		}
	}
}
