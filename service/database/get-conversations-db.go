package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetConversations(userId int) ([]ConversationDB, string, error) {
	var response []ConversationDB

	query := `
		SELECT c.conversation_id, c.type, c.message_id, c.group_name, c.group_image
		FROM conversations c
		INNER JOIN participants p ON c.conversation_id = p.conversation_id
		WHERE p.user_id = ?;`

	rows, err := db.c.Query(query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "no conversations found", ErrNotFound
		}
		return nil, "internal error while executing the conversations search query", ErrInternal
	}

	for rows.Next() {
		var conv ConversationDB
		var mess MessageDB
		var messageID sql.NullInt64
		var comm []CommentDB
		var chatName sql.NullString
		var chatImageUrl sql.NullString

		if err := rows.Scan(&conv.ConvID, &conv.ChatType, &messageID, &chatName, &chatImageUrl); err != nil {
			_ = rows.Close()
			return nil, "internal error while scanning conversation data " + err.Error(), ErrInternal
		}

		if messageID.Valid {
			mess.MessageID = int(messageID.Int64)
		}

		conv.ChatName = chatName.String
		conv.ChatImageUrl = chatImageUrl.String

		// --- Recupero membri ---
		membersQuery := `
			SELECT u.username
			FROM participants p
			INNER JOIN users u ON p.user_id = u.user_id
			WHERE p.conversation_id = ?;`

		memberRows, err := db.c.Query(membersQuery, conv.ConvID)
		if err != nil {
			_ = rows.Close()
			return nil, "internal error while fetching conversation members", ErrInternal
		}

		for memberRows.Next() {
			var member string
			if err := memberRows.Scan(&member); err != nil {
				_ = memberRows.Close()
				_ = rows.Close()
				return nil, "internal error while scanning members", ErrInternal
			}
			conv.Members = append(conv.Members, member)
		}
		if err := memberRows.Err(); err != nil {
			_ = memberRows.Close()
			_ = rows.Close()
			return nil, "errore durante la scan delle righe " + err.Error(), ErrInternal
		}
		_ = memberRows.Close()

		if conv.ChatType == "private" {
			query := `
				SELECT u.username, u.url_profile_image
				FROM participants p
				INNER JOIN users u ON p.user_id = u.user_id
				WHERE p.conversation_id = ? AND p.user_id != ?;`

			userRow := db.c.QueryRow(query, conv.ConvID, userId)
			if err := userRow.Scan(&conv.ChatName, &conv.ChatImageUrl); err != nil {
				_ = rows.Close()
				return nil, "internal error while fetching user details", ErrInternal
			}
		}

		// --- Recupero messaggio ---
		if mess.MessageID > 0 {
			query = `
				SELECT u.username, m.type, m.timestamp, m.content_text, m.content_image, m.is_forwarded, m.reply_message_id
				FROM messages m
				INNER JOIN users u ON m.user_id = u.user_id
				WHERE m.message_id = ?;`

			row := db.c.QueryRow(query, mess.MessageID)
			var text sql.NullString
			var imageUrl sql.NullString
			var isForwarded bool
			var replyTo sql.NullInt64

			if err = row.Scan(&mess.SenderUsername, &mess.MessageType, &mess.Timestamp, &text, &imageUrl, &isForwarded, &replyTo); err != nil {
				_ = rows.Close()
				return nil, "internal error while scanning message data: " + err.Error(), ErrInternal
			}

			mess.Text = text.String
			mess.ImageUrl = imageUrl.String
			mess.IsForwarded = isForwarded

			if replyTo.Valid {
				mess.IsAnswering = int(replyTo.Int64)
			}

			readQuery := `
				SELECT is_delivered, is_read
				FROM messageReads
				WHERE message_id = ? AND user_id = ?;`

			err = db.c.QueryRow(readQuery, mess.MessageID, userId).Scan(&mess.IsDelivered, &mess.IsRead)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				_ = rows.Close()
				return nil, "errore nel recupero di is_read/is_delivered: " + err.Error(), ErrInternal
			}

			// --- Reactions ---
			query = `
				SELECT r.user_id, r.unicode
				FROM reactions r
				WHERE r.message_id = ?;`

			commRows, err := db.c.Query(query, mess.MessageID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				_ = rows.Close()
				return nil, "internal error while retrieving reactions", ErrInternal
			}
			if commRows != nil {
				for commRows.Next() {
					var singleComm CommentDB
					var userID int
					if err = commRows.Scan(&userID, &singleComm.Emoji); err != nil {
						_ = commRows.Close()
						_ = rows.Close()
						return nil, "internal error while scanning message data", ErrInternal
					}
					singleComm.Username, _ = db.GetUsernameFromUid(userID)
					comm = append(comm, singleComm)
				}
				if err := commRows.Err(); err != nil {
					_ = commRows.Close()
					_ = rows.Close()
					return nil, "errore durante la scan delle reazioni " + err.Error(), ErrInternal
				}
				_ = commRows.Close()
			}

			mess.Comments = comm
			if len(mess.Comments) == 0 {
				mess.Comments = []CommentDB{}
			}
		}

		conv.LastMessage = mess
		unreadQuery := `
			SELECT COUNT(*)
			FROM messageReads mr
			INNER JOIN messages m ON mr.message_id = m.message_id
			WHERE m.conversation_id = ? AND mr.user_id = ? AND mr.is_read = FALSE;`

		err = db.c.QueryRow(unreadQuery, conv.ConvID, userId).Scan(&conv.UnreadMessages)
		if err != nil {
			_ = rows.Close()
			return nil, "errore nel calcolo dei messaggi non letti: " + err.Error(), ErrInternal
		}

		response = append(response, conv)
	}

	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, "errore durante la scansione delle conversazioni: " + err.Error(), ErrInternal
	}
	_ = rows.Close()

	/*
		unreadQuery := `
			UPDATE messageReads
			SET is_delivered = TRUE
			WHERE user_id = ? AND is_delivered = FALSE;
			`

		_, err = db.c.Exec(unreadQuery, userId)
		if err != nil {
			return nil, "errore nell'impostazione di ricezione dei messaggi: " + err.Error(), ErrInternal
		}

	*/

	return response, "conversations correctly returned", nil
}
