package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
	"net/http"
	"strconv"
)

const (
	DefaultMessageLimit = 200
)

// POST /channels/:channelId/messages メッセージをチャンネルに投稿、embedを自動変換（embed=1）
func PostMessage(channelId, content string) (*openapi.Message, error) {
	message, _, err := client.MessageApi.PostMessage(
		auth,
		channelId,
		&openapi.MessageApiPostMessageOpts{
			PostMessageRequest: optional.NewInterface(openapi.PostMessageRequest{
				Content: content,
				Embed:   true,
			}),
		},
	)
	return &message, err
}

// PUT /messages/:messageId メッセージを編集
func EditMessage(messageId, content string) (*openapi.Message, error) {
	message, _, err := client.MessageApi.EditMessage(
		auth,
		messageId,
		&openapi.MessageApiEditMessageOpts{
			PostMessageRequest: optional.NewInterface(openapi.PostMessageRequest{
				Content: content,
				Embed:   true,
			}),
		},
	)
	return &message, err
}

// GET /channels/:channelId/messages チャンネルのメッセージを取得 order: desc, 時間が新しい方から古い方へ
func GetChannelMessages(channelId string, limit, offset int) (messages []openapi.Message, hasMore bool, err error) {
	var res *http.Response
	messages, res, err = client.MessageApi.GetMessages(auth, channelId, &openapi.MessageApiGetMessagesOpts{
		Limit:  optional.NewInt32(int32(limit)),
		Offset: optional.NewInt32(int32(offset)),
		Order:  optional.NewString("desc"),
	})
	if err != nil {
		return
	}

	hasMore, err = strconv.ParseBool(res.Header.Get("x-traq-more"))
	return
}
