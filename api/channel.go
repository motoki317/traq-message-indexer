package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

var (
	channelsCache *ChannelsMap
)

type ChannelsMap struct {
	Public map[string]*openapi.Channel
	Dm     map[string]*openapi.DmChannel
}

// GET /channels チャンネルリストを取得
func GetChannels(includeDMs bool, canUseCache bool) (*ChannelsMap, error) {
	if canUseCache && channelsCache != nil {
		return channelsCache, nil
	}

	channels, _, err := client.ChannelApi.GetChannels(
		auth,
		&openapi.ChannelApiGetChannelsOpts{
			IncludeDm: optional.NewBool(includeDMs),
		},
	)

	channelsMap := makeChannelsMap(&channels)
	channelsCache = channelsMap
	return channelsMap, err
}

func makeChannelsMap(channelList *openapi.ChannelList) *ChannelsMap {
	ret := ChannelsMap{
		Public: make(map[string]*openapi.Channel),
		Dm:     make(map[string]*openapi.DmChannel),
	}

	for _, ch := range channelList.Public {
		channel := ch
		ret.Public[channel.Id] = &channel
	}
	for _, ch := range channelList.Dm {
		channel := ch
		ret.Dm[channel.Id] = &channel
	}

	return &ret
}
