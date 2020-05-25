package handler

import (
	"fmt"
	"github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const messagesPerPage = 5

// extractKeywords scans arguments and extracts valid keywords.
func extractKeywords(args []string) []string {
	keywords := make([]string, 0)

	skipNext := false
	for _, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if arg == "-p" || arg == "--page" {
			// page
			skipNext = true
			continue
		} else if strings.HasPrefix(arg, "-") {
			// special args
			continue
		} else if arg == "this" || strings.HasPrefix(arg, "#") {
			// channel
			continue
		} else if strings.HasPrefix(arg, "@") || strings.HasPrefix(arg, "from:") {
			// user
			continue
		}

		keywords = append(keywords, arg)
	}

	return keywords
}

// extractPage scans arguments (expects to be written in 1-indexed number), and returns 0-indexed page number if exists.
func extractPage(args []string) int {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-p" || args[i] == "--page" {
			if page, err := strconv.Atoi(args[i+1]); err == nil {
				// from 1-indexed page to 0-indexed page
				return page - 1
			}
		}
	}
	return 0
}

// extractChannelIDs extracts mentioned channels from the message payload, and returns channel ids slice.
func extractChannelIDs(payload *traqbot.MessageCreatedPayload, args []string) []string {
	includesThis := false
	for _, arg := range args {
		if strings.ToLower(arg) == "this" {
			includesThis = true
		}
	}

	channelIDs := make([]string, 0)
	if includesThis {
		channelIDs = append(channelIDs, payload.Message.ChannelID)
	}

	// Check embedded channel links
	for _, e := range payload.Message.Embedded {
		if e.Type == "channel" {
			channelIDs = append(channelIDs, e.ID)
		}
	}

	return channelIDs
}

// extractUserIDs extract user IDs from the message.
func extractUserIDs(payload *traqbot.MessageCreatedPayload, args []string) []string {
	userIDs := make([]string, 0)

	for _, e := range payload.Message.Embedded {
		if e.Type == "user" {
			userIDs = append(userIDs, e.ID)
		}
	}

	from := make([]string, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, "from:") {
			from = append(from, arg[len("from:"):])
		}
	}

	// If specified by user name, request the API
	if len(from) > 0 {
		users, err := api.GetNameUserMap(true, true)
		if err == nil {
			for _, name := range from {
				if user, ok := users[name]; ok {
					userIDs = append(userIDs, user.Id)
				}
			}
		} else {
			log.Printf("an error occurred while requesting users: %s\n", err)
		}
	}

	return userIDs
}

type format struct {
	re     *regexp.Regexp
	layout string
}

var formats = []format{
	{
		re:     regexp.MustCompile("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{2}:\\d{2}"),
		layout: "2006/01/02 15:04:05 +0900",
	},
	{
		re:     regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{1,2}:\\d{2}:\\d{2}"),
		layout: "2006-01-02 15:04:05 +0900",
	},
	{
		re:     regexp.MustCompile("\\d{4}/\\d{2}/\\d{2} \\d{1,2}:\\d{2}"),
		layout: "2006/01/02 15:04 +0900",
	},
	{
		re:     regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{1,2}:\\d{2}"),
		layout: "2006-01-02 15:04 +0900",
	},
	{
		re:     regexp.MustCompile("\\d{4}/\\d{2}/\\d{2}"),
		layout: "2006/01/02 +0900",
	},
	{
		re:     regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}"),
		layout: "2006-01-02 +0900",
	},
}

func parseTime(s string) (t time.Time, ok bool) {
	for _, f := range formats {
		if matched := f.re.FindString(s); matched != "" {
			if t, err := time.Parse(f.layout, matched+" +0900"); err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

func extractTimeRange(args []string) (after, before *time.Time) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "after:") {
			if t, ok := parseTime(arg[len("after:"):]); ok {
				t = t.In(time.UTC)
				after = &t
			}
		} else if strings.HasPrefix(arg, "before:") {
			if t, ok := parseTime(arg[len("before:"):]); ok {
				t = t.In(time.UTC)
				before = &t
			}
		}
	}
	return
}

func commandSearch() command {
	return command{
		name: "search",
		help: strings.Join([]string{
			"## Search Command Help",
			"Searches through messages.",
			"Syntax: `/search <keywords> [#channel filter] [-c|--children] [-r|--recursive] [-p|--page #]`",
			"",
			"### Specify keywords",
			"Specify keywords using the following syntax ('IN BOOLEAN MODE'): https://mariadb.com/kb/en/full-text-index-overview/#in-boolean-mode",
			"The search is case-insensitive.",
			"- /search じゃんけん",
			"",
			"### Filter by channels",
			"Filter by channels by mentioning channels. Append '--children' (and '--recursive') argument(s) to instead specify its child channels.",
			"- /search 講習会 #general",
			"- /search おいす #gps/times --children --recursive",
			"    - Recursively searches channels under #gps/times, but NOT including #gps/times.",
			"",
			"### Filter by users",
			"Filter by users by mentioning them, or using syntax `from:<user id>`.",
			"- /search おいす @toki",
			"- /search おいす from:kashiwade from:liquid1224",
			"",
			"### Filter by time",
			"Filter by time using syntax `after:<time>` and/or `before:<time>`.",
			"- /search おいす before:2020/01/01",
			"- /search 講習会 after:2020/04/01 07:00:00",
		}, "\n"),
		handle: func(h *handler, payload *traqbot.MessageCreatedPayload, args []string) error {
			keywords := extractKeywords(args[1:])
			page := extractPage(args[1:])
			channelIDs := extractChannelIDs(payload, args[1:])
			userIDs := extractUserIDs(payload, args[1:])
			after, before := extractTimeRange(args[1:])

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

			count, err := h.repo.SearchMessageCount(keywords, channelIDs, userIDs, after, before)
			if err != nil {
				return err
			}
			messages, err := h.repo.SearchMessage(keywords, channelIDs, userIDs, after, before, messagesPerPage, page*messagesPerPage)
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
