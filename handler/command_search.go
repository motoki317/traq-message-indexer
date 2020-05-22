package handler

import (
	"fmt"
	"github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	"log"
	"strconv"
	"strings"
	"time"
)

const messagesPerPage = 5

func commandSearch() command {
	return command{
		name: "search",
		help: "Searches through messages.\n" +
			"Syntax: `/search <keywords> [#channel filter] [-c|--children] [-r|--recursive] [-p|--page #]`\n" +
			"Examples:\n" +
			"- `/search じゃんけん`\n" +
			"- `/search じゃんけん #gps/times/kashiwade`\n" +
			"- `/search おいすー #gps/times --children --recursive --page 2`",
		handle: func(h *handler, payload *traqbot.MessageCreatedPayload, args []string) error {
			keywords := make([]string, 0)
			page := 0
			// parse args to extract keywords
			for i := 1; i < len(args); i++ {
				if args[i] == "-p" || args[i] == "--page" {
					i++
					// parse page
					if i < len(args) {
						if parsed, err := strconv.Atoi(args[i]); err == nil {
							// to 0-indexed
							page = parsed - 1
						}
					}
					continue
				} else if strings.HasPrefix(args[i], "-") {
					continue
				} else if args[i] == "this" || strings.HasPrefix(args[i], "#") {
					continue
				}
				keywords = append(keywords, args[i])
			}

			channelIDs := extractMentionedChannels(payload)

			channels, err := api.GetChannels(false, true)
			if err != nil {
				return err
			}

			// Parse args and extract channel IDs
			argMap := makeArgMap(args)
			getChild := argMap["-c"] || argMap["--children"]
			recursive := argMap["-r"] || argMap["--recursive"]
			if getChild {
				childIDs := make([]string, 0)
				for _, id := range channelIDs {
					childIDs = append(childIDs, getChildChannels(id, channels, recursive)...)
				}
				channelIDs = childIDs
			}

			if len(keywords) == 0 {
				if _, err := api.PostMessage(payload.Message.ChannelID, "Please specify at least one keyword!"); err != nil {
					log.Println(err)
				}
				return nil
			}

			count, err := h.repo.SearchMessageCount(keywords, channelIDs)
			if err != nil {
				return err
			}
			messages, err := h.repo.SearchMessage(keywords, channelIDs, messagesPerPage, page*messagesPerPage)
			if err != nil {
				return err
			}
			if len(messages) == 0 {
				if _, err := api.PostMessage(payload.Message.ChannelID, "No messages found."); err != nil {
					log.Println(err)
				}
				return nil
			}

			// Format final response
			/**
			**Search results**

			- #gps/times/kashiwade @ 2020-05-08 11:48:24.328488Z
			https://q.trap.jp/messages/feeae600-b2ca-46af-973c-2482166d7b67

			*Page 1 of 5, showing 1 ~ 5 of 23 results*
			*/
			maxPage := (count - 1) / messagesPerPage
			formatted := make([]string, 0, len(messages))
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			for _, message := range messages {
				var channelPath string
				if channel, ok := channels.Public[message.ChannelID]; ok {
					channelPath = getChannelPath(channel, channels)
				} else {
					channelPath = "#???"
				}

				formatted = append(formatted, fmt.Sprintf(
					"- %s @ %s\n"+
						"https://q.trap.jp/messages/%s",
					channelPath, message.CreatedAt.In(jst).Format("2006-01-02 15:04:05"),
					message.ID,
				))
			}

			finalMessage := fmt.Sprintf(
				"**Search results**\n"+
					"\n"+
					"%s\n"+
					"\n"+
					"*Page %v of %v, showing %v ~ %v of %v results*",
				strings.Join(formatted, "\n"),
				page+1, maxPage+1, page*messagesPerPage+1, page*messagesPerPage+len(messages),
				count,
			)

			if _, err := api.PostMessage(payload.Message.ChannelID, finalMessage); err != nil {
				log.Println(err)
			}
			return nil
		},
	}
}
