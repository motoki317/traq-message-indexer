package handler

import (
	"github.com/motoki317/traq-message-indexer/api"
	"github.com/motoki317/traq-message-indexer/repository"
	"log"
	"time"
)

// channelIDのメッセージ全てを更新する
func (h *handler) updateChannel(channelID string) error {
	h.repo.ChannelLock()
	defer h.repo.ChannelUnlock()

	seen, err := h.repo.GetSeenChannel(channelID)
	if err != nil {
		return err
	}
	if seen == nil {
		seen = &repository.SeenChannel{
			ID:                   channelID,
			LastProcessedMessage: time.Unix(0, 0),
		}
	}

	// LastProcessedMessageかそれより前に送られたメッセージににたどり着くか、
	// チャンネルの最後まで読み込む
	processUntil := time.Now()
	limit := api.DefaultMessageLimit
	offset := 0
	for {
		log.Printf("Processing for channel %s, limit %v, offset %v...\n", channelID, limit, offset)
		messages, hasMore, err := api.GetChannelMessages(channelID, limit, offset)
		if err != nil {
			return err
		}

		processedAll := false
		for _, m := range messages {
			if m.CreatedAt.After(processUntil) {
				continue
			}
			if m.CreatedAt.Equal(seen.LastProcessedMessage) || m.CreatedAt.Before(seen.LastProcessedMessage) {
				processedAll = true
				break
			}

			message := m
			err := h.repo.CreateMessage(&repository.Message{
				ID:        message.Id,
				ChannelID: message.ChannelId,
				UserID:    message.UserId,
				CreatedAt: message.CreatedAt,
				Text:      message.Content,
			})
			if err != nil {
				return err
			}
		}

		if processedAll || !hasMore {
			break
		}

		offset += limit

		// 1秒間隔でメッセージを取りに行く
		<-time.NewTimer(time.Second * 1).C
	}

	log.Printf("Processed for channel %s, from %v to %v", channelID, seen.LastProcessedMessage, processUntil)

	// どこまで見たかを保存する
	seen.LastProcessedMessage = processUntil
	return h.repo.UpdateSeenChannel(seen)
}
