package handler

import (
	"fmt"
	"github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	"log"
	"strings"
	"sync"
	"time"
)

func commandIndex() command {
	return command{
		name: "index",
		help: "Indexes a channel's all messages.\n" +
			"Examples:\n" +
			"- `/index this`\n" +
			"- `/index #gps/times/Yataka_ML`\n" +
			"- `/index #gps/times --children`\n" +
			"- `/index #gps/times --children --recursive`",
		handle: func(h *handler, payload *traqbot.MessageCreatedPayload, args []string) error {
			channelIDs := extractMentionedChannels(payload)

			channels, err := api.GetChannels(false, false)
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

			// Send message
			toEdit, err := api.PostMessage(payload.Message.ChannelID, fmt.Sprintf(
				"Indexing %v channels... (%v/%v) :rotarex:", len(channelIDs), 0, len(channelIDs)))
			if err != nil {
				return err
			}

			editMut := sync.Mutex{}
			ticker := time.NewTicker(time.Second * 5)
			// true if finished correctly, false if not
			finish := make(chan bool)
			done := 0

			// Update each channel
			go func() {
				for _, id := range channelIDs {
					err := h.updateChannel(id)
					// Error on updating
					if err != nil {
						editMut.Lock()
						finish <- false
						if _, err := api.PostMessage(payload.Message.ChannelID, fmt.Sprintf(
							"Failed to update a channel's messages... channel id: %s\n%s", id, err)); err != nil {
							log.Println(err)
						}
						editMut.Unlock()
						return
					}
					done++
				}
				finish <- true
			}()

			go func() {
				for {
					select {
					case <-ticker.C:
						editMut.Lock()
						if _, err := api.EditMessage(toEdit.Id, fmt.Sprintf(
							"Indexing %v channels... (%v/%v) :rotarex:", len(channelIDs), done, len(channelIDs))); err != nil {
							log.Println(err)
						}
						editMut.Unlock()
					case done := <-finish:
						if done {
							if _, err := api.EditMessage(toEdit.Id, fmt.Sprintf(
								"Finished indexing %v channels!\n%s", len(channelIDs), formatChannelIDs(channelIDs, channels))); err != nil {
								log.Println(err)
							}
						} else {
							if _, err := api.EditMessage(toEdit.Id, "An error occurred..."); err != nil {
								log.Println(err)
							}
						}

						ticker.Stop()
						return
					}
				}
			}()

			return nil
		},
	}
}

func formatChannelIDs(channelIDs []string, channels *api.ChannelsMap) string {
	formatted := make([]string, 0, len(channelIDs))
	for _, id := range channelIDs {
		ch, ok := channels.Public[id]
		if !ok {
			continue
		}
		formatted = append(formatted, "`"+getChannelPath(ch, channels)+"`")
	}

	return strings.Join(formatted, ", ")
}
