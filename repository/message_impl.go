package repository

import (
	"database/sql"
	"strings"
	"time"
)

func (r *RepositoryImpl) CreateMessage(message *Message) error {
	_, err := r.db.Exec("INSERT INTO `message` (id, channel_id, user_id, created_at, text) VALUES (?, ?, ?, ?, ?)",
		message.ID, message.ChannelID, message.UserID, message.CreatedAt, message.Text)
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

func nCopies(s string, n int) []string {
	ret := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ret = append(ret, s)
	}
	return ret
}

// buildSearchMessageQuery builds a partial query ('WHERE' clause) for searching message.
func (r *RepositoryImpl) buildSearchMessageQuery(keywords, channelIDs, userIDs []string, after, before *time.Time) (string, []interface{}) {
	ands := make([]string, 0)
	args := make([]interface{}, 0)

	if len(keywords) > 0 {
		ands = append(ands, "MATCH (`text`) AGAINST (? IN BOOLEAN MODE)")
		args = append(args, strings.Join(keywords, " "))
	}

	if len(channelIDs) > 0 {
		// has one of the specified channel IDs
		ands = append(ands, "("+strings.Join(nCopies("`channel_id` = ?", len(channelIDs)), " OR ")+")")
		for _, channelID := range channelIDs {
			args = append(args, channelID)
		}
	}

	if len(userIDs) > 0 {
		// has one of the specified user IDs
		ands = append(ands, "("+strings.Join(nCopies("`user_id` = ?", len(userIDs)), " OR ")+")")
		for _, userID := range userIDs {
			args = append(args, userID)
		}
	}

	if after != nil {
		ands = append(ands, "`created_at` > ?")
		args = append(args, after.Format("2006-01-02 15:04:05"))
	}
	if before != nil {
		ands = append(ands, "`created_at` < ?")
		args = append(args, before.Format("2006-01-02 15:04:05"))
	}

	// No 'WHERE' clause
	if len(ands) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(ands, " AND "), args
}

func (r *RepositoryImpl) SearchMessage(keywords, channelIDs, userIDs []string, after, before *time.Time, limit, offset int) ([]Message, error) {
	where, args := r.buildSearchMessageQuery(keywords, channelIDs, userIDs, after, before)

	args = append(args, limit)
	args = append(args, offset)

	var messages []Message
	if err := r.db.Select(&messages, "SELECT * FROM `message` "+where+" ORDER BY `created_at` DESC LIMIT ? OFFSET ?", args...); err != nil && err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return messages, nil
	}
}

func (r *RepositoryImpl) SearchMessageCount(keywords, channelIDs, userIDs []string, after, before *time.Time) (int, error) {
	where, args := r.buildSearchMessageQuery(keywords, channelIDs, userIDs, after, before)

	var count int
	if err := r.db.Get(&count, "SELECT COUNT(*) FROM `message` "+where, args...); err != nil && err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	} else {
		return count, nil
	}
}
