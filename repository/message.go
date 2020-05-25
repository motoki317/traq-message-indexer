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
	// メッセージを検索します。
	SearchMessage(keywords, channelIDs, userIDs []string, after, before *time.Time, limit, offset int) ([]Message, error)
	// 検索結果のメッセージの総数を返します。
	SearchMessageCount(keywords, channelIDs, userIDs []string, after, before *time.Time) (int, error)
}
