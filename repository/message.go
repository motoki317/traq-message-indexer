package repository

import "time"

type Message struct {
	ID        string    `db:"id"`
	ChannelID string    `db:"channel_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	Text      string    `db:"text"`
}

type MessageRepository interface {
	// メッセージを作成しストアします。
	CreateMessage(*Message) error
	// メッセージをIDで取得します。
	GetMessage(messageID string) (*Message, error)
	// キーワードと指定されたチャンネルIDからメッセージを検索します。
	SearchMessage(keywords, channelIDs, userIDs []string, limit, offset int) ([]Message, error)
	// キーワードと指定されたチャンネルIDから検索されるメッセージの総数を返します。
	SearchMessageCount(keywords, channelIDs, userIDs []string) (int, error)
}
