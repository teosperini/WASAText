package database

import (
	"database/sql"
	"errors"
	"strings"
)

func (db *appdbimpl) ValidateID(ID int) error {
	_, err := db.GetUsernameFromUid(ID)
	return err
}

func (db *appdbimpl) GetUidFromUsername(username string) (int, error) {
	var uid int
	query := "SELECT user_id FROM users WHERE username = ?"
	err := db.c.QueryRow(query, username).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		} else {
			return 0, ErrInternal
		}
	}
	return uid, nil
}

func (db *appdbimpl) GetUsernameFromUid(id int) (string, error) {
	var name string
	query := "SELECT username FROM users WHERE user_id = ?"
	err := db.c.QueryRow(query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", ErrInternal
		}
	}
	return name, nil
}

func (db *appdbimpl) GetUrlFromUid(id int) (string, error) {
	var url string
	query := "SELECT url_profile_image FROM users WHERE user_id = ?"
	err := db.c.QueryRow(query, id).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", ErrInternal
		}
	}
	return url, nil
}
func (db *appdbimpl) IsUniqueConstraintError(err error) bool {
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return true
	}
	return false
}

func (db *appdbimpl) IsUserInTheConversation(convId int, userId int) (string, error) {
	// Controllo se l'utente appartiene alla conversazione
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM participants
			WHERE conversation_id = ? AND user_id = ?
		)`
	err := db.c.QueryRow(query, convId, userId).Scan(&exists)
	if err != nil {
		return "error while searching for the user in the conversation " + err.Error(), ErrInternal
	}
	if !exists {
		return "utente non autorizzato", ErrForbidden
	}

	return "", nil
}

func (db *appdbimpl) MarkMessagesAsDelivered(userId int, conversations []ConversationDB) error {
	for _, conv := range conversations {
		if !conv.LastMessage.Timestamp.Valid {
			continue
		}
		username, err := db.GetUsernameFromUid(userId)
		if err != nil {
			return err
		}
		// Non marcare come "delivered" i miei messaggi
		if conv.LastMessage.SenderUsername == "" || conv.LastMessage.SenderUsername == username {
			continue
		}

		query := `
			UPDATE messageReads
			SET is_delivered = TRUE
			WHERE user_id = ?
			  AND message_id IN (
				SELECT m.message_id
				FROM messages m
				WHERE m.conversation_id = ?
				  AND m.timestamp <= ?
				  AND m.user_id != ?
			  )
			  AND is_delivered = FALSE;
		`

		_, err = db.c.Exec(query, userId, conv.ConvID, conv.LastMessage.Timestamp.Time, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *appdbimpl) GetMessageById(messageId int) (MessageDB, error) {
	query := `
	SELECT m.message_id, u.username, m.timestamp, m.type, m.content_text, m.content_image
	FROM messages m
	INNER JOIN users u ON m.user_id = u.user_id
	WHERE m.message_id = ?`

	row := db.c.QueryRow(query, messageId)

	var msg MessageDB
	var text sql.NullString
	var imageUrl sql.NullString

	if err := row.Scan(&msg.MessageID, &msg.SenderUsername, &msg.Timestamp, &msg.MessageType, &text, &imageUrl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MessageDB{}, ErrNotFound
		}
		return MessageDB{}, err
	}

	msg.Text = text.String
	msg.ImageUrl = imageUrl.String

	return msg, nil
}
