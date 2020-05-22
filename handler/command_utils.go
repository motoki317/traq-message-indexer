package handler

import (
	traqbot "github.com/motoki317/traq-bot"
	"github.com/motoki317/traq-message-indexer/api"
	traqApi "github.com/sapphi-red/go-traq"
	"strings"
)

// makeArgMap returns map of arguments with each string key in lower case.
func makeArgMap(args []string) (ret map[string]bool) {
	ret = make(map[string]bool)
	for _, v := range args {
		ret[strings.ToLower(v)] = true
	}
	return
}

// extractMentionedChannel extracts mentioned channels from the message payload, and returns channel ids slice.
func extractMentionedChannels(payload *traqbot.MessageCreatedPayload) []string {
	// Check "this" argument"
	fields := strings.Fields(payload.Message.PlainText)
	includesThis := false
	for _, f := range fields {
		if strings.ToLower(f) == "this" {
			includesThis = true
		}
	}

	ret := make([]string, 0)
	if includesThis {
		ret = append(ret, payload.Message.ChannelID)
	}

	// Check embedded channel links
	for _, e := range payload.Message.Embedded {
		if e.Type == "channel" {
			ret = append(ret, e.ID)
		}
	}

	return ret
}

// getChildChannels returns a list of child channel IDs. Does not include the parent.
func getChildChannels(channelId string, channels *api.ChannelsMap, recursive bool) []string {
	res := make([]string, 0)

	ch, ok := channels.Public[channelId]
	if !ok {
		return res
	}

	res = append(res, ch.Children...)
	if recursive {
		for _, child := range ch.Children {
			res = append(res, getChildChannels(child, channels, recursive)...)
		}
	}

	return res
}

// getChannelPath returns a channel path of the given channel in string.
func getChannelPath(channel *traqApi.Channel, channels *api.ChannelsMap) string {
	if channel.ParentId == nil {
		return "#" + channel.Name
	}

	parent, ok := channels.Public[*channel.ParentId]
	if !ok {
		return "unknown channel"
	}

	return getChannelPath(parent, channels) + "/" + channel.Name
}
