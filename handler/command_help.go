package handler

import (
	"github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	"log"
	"strings"
)

func formatCommands(commands []command) string {
	formatted := make([]string, 0, len(commands))
	for _, c := range commands {
		formatted = append(formatted, "`"+prefix+c.name+"`")
	}
	return strings.Join(formatted, ", ")
}

func commandHelp(commands []command) command {
	help := "Shows help of each command.\n" +
		"Syntax: `/help <command name>`\n" +
		"Examples: `/help index`, `/help help`\n" +
		"Commands list: " + formatCommands(commands)

	return command{
		name: "help",
		help: help,
		handle: func(h *handler, payload *traqbot.MessageCreatedPayload, args []string) error {
			channelID := payload.Message.ChannelID
			if len(args) <= 1 {
				if _, err := api.PostMessage(channelID, help); err != nil {
					log.Println(err)
				}
				return nil
			}

			specified := strings.ToLower(args[1])
			for _, c := range commands {
				if c.name == specified {
					if _, err := api.PostMessage(channelID, c.help); err != nil {
						log.Println(err)
					}
					return nil
				}
			}

			if _, err := api.PostMessage(channelID, "Unknown command."); err != nil {
				log.Println(err)
			}
			return nil
		},
	}
}
