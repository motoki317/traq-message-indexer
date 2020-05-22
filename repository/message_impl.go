package repository

import (
	"database/sql"
	"strings"
)

func (r *RepositoryImpl) CreateMessage(message *Message) error {
	_, err := r.db.Exec("INSERT INTO `message` (id, channel_id, created_at, text) VALUES (?, ?, ?, ?)",
		message.ID, message.ChannelID, message.CreatedAt, message.Text)
	return err
}

func (r *RepositoryImpl) GetMessage(messageID string) (*Message, error) {
	var message Message
	if err := r.db.Get(&message, "SELECT * FROM `message` WHERE `id` = ?", messageID); err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &message, nil
	}
}

func (r *RepositoryImpl) SearchMessage(keywords []string, channelIDs []string, limit, offset int) ([]Message, error) {
	if len(keywords) == 0 {
		return nil, nil
	}

	for i, keyword := range keywords {
		keywords[i] = "+" + keyword
	}
	matchAgainst := strings.Join(keywords, " ")

	channelIDRestrictions := strings.Repeat(" AND `channel_id` = ?", len(channelIDs))

	args := []interface{}{matchAgainst}
	for _, id := range channelIDs {
		args = append(args, id)
	}
	args = append(args, limit)
	args = append(args, offset)

	var messages []Message
	if err := r.db.Select(&messages, "SELECT * FROM `message` WHERE MATCH (`text`) AGAINST (? IN BOOLEAN MODE)"+channelIDRestrictions+" ORDER BY `created_at` DESC LIMIT ? OFFSET ?", args...); err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return messages, nil
	}
}

func (r *RepositoryImpl) SearchMessageCount(keywords []string, channelIDs []string) (int, error) {
	if len(keywords) == 0 {
		return 0, nil
	}

	for i, keyword := range keywords {
		keywords[i] = "+" + keyword
	}
	matchAgainst := strings.Join(keywords, " ")

	channelIDRestrictions := strings.Repeat(" AND `channel_id` = ?", len(channelIDs))

	args := []interface{}{matchAgainst}
	for _, id := range channelIDs {
		args = append(args, id)
	}

	var count int
	if err := r.db.Select(&count, "SELECT * FROM `message` WHERE MATCH (`text`) AGAINST (? IN BOOLEAN MODE)"+channelIDRestrictions, args...); err != nil && err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	} else {
		return count, nil
	}
}
