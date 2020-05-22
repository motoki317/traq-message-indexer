package repository

import "time"

type Message struct {
	ID        string    `db:"id"`
	ChannelID string    `db:"channel_id"`
	CreatedAt time.Time `db:"created_at"`
	Text      string    `db:"text"`
}

type MessageRepository interface {
	// メッセージを作成しストアします。
	CreateMessage(*Message) error
	// メッセージをIDで取得します。
	GetMessage(messageID string) (*Message, error)
	// キーワードと指定されたチャンネルIDからメッセージを検索します。
	SearchMessage(keywords []string, channelIDs []string, limit, offset int) ([]Message, error)
	// キーワードと指定されたチャンネルIDから検索されるメッセージの総数を返します。
	SearchMessageCount(keywords []string, channelIDs []string) (int, error)
}
