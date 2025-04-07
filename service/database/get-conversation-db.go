package database

import (
	"database/sql"
	"sort"
)

func (db *appdbimpl) GetConversationDB(userId int, convId int) ([]MessageDB, string, error) {
	// 1. Verifica se l'utente fa parte della conversazione
	if _, err := db.IsUserInTheConversation(convId, userId); err != nil {
		return nil, "utente non parte della conversazione", err
	}

	// 2. Prendi il tipo di conversazione (private/group)
	var chatType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE conversation_id = ?`, convId).Scan(&chatType)
	if err != nil {
		return nil, "errore nel determinare il tipo di chat: " + err.Error(), ErrInternal
	}

	// 3. Aggiorna i messaggi ricevuti (da altri utenti) come letti per userId
	_, err = db.c.Exec(`
		UPDATE messageReads
		SET is_read = TRUE
		WHERE message_id IN (
			SELECT message_id FROM messages
			WHERE conversation_id = ? AND user_id != ?
		) AND user_id = ?`, convId, userId, userId)
	if err != nil {
		return nil, "errore nell'aggiornamento dello stato di lettura: " + err.Error(), ErrInternal
	}

	// 4. Ottieni il tuo username (ti serve per capire se sei il mittente del messaggio)
	username, err := db.GetUsernameFromUid(userId)
	if err != nil {
		return nil, "errore nel recupero dell'username", err
	}

	// 5. Prepara strutture dati
	var messages []MessageDB
	messageMap := make(map[int]*MessageDB)

	// 6. Query messaggi + reazioni
	rows, err := db.c.Query(`
		SELECT
			m.message_id, u_sender.username, m.timestamp, m.type, m.content_text, m.content_image,
			u_react.username, r.unicode, m.is_forwarded, m.reply_message_id
		FROM messages m
		JOIN users u_sender ON m.user_id = u_sender.user_id
		LEFT JOIN reactions r ON r.message_id = m.message_id
		LEFT JOIN users u_react ON r.user_id = u_react.user_id
		WHERE m.conversation_id = ?
		ORDER BY m.timestamp`, convId)
	if err != nil {
		return nil, "errore nel recupero dei messaggi: " + err.Error(), ErrInternal
	}

	for rows.Next() {
		var (
			msgID                     sql.NullInt64
			senderUsername, msgType   sql.NullString
			contentText, contentImage sql.NullString
			reactionUsername, emoji   sql.NullString
			timestamp                 sql.NullTime
		)

		var isForwarded bool
		var replyTo sql.NullInt64

		if err := rows.Scan(&msgID, &senderUsername, &timestamp, &msgType, &contentText, &contentImage, &reactionUsername, &emoji, &isForwarded, &replyTo); err != nil {
			return nil, "errore nella scansione dei messaggi: " + err.Error(), ErrInternal
		}

		id := int(msgID.Int64)

		// Se il messaggio non esiste ancora, lo inizializziamo
		if _, exists := messageMap[id]; !exists {
			messageMap[id] = &MessageDB{
				MessageID:      id,
				SenderUsername: senderUsername.String,
				Timestamp:      timestamp,
				MessageType:    msgType.String,
				Text:           contentText.String,
				ImageUrl:       contentImage.String,
				Comments:       []CommentDB{},
			}
		}

		msg := messageMap[id]
		msg.IsForwarded = isForwarded
		if replyTo.Valid {
			msg.IsAnswering = int(replyTo.Int64)
		}

		// Reazioni
		if reactionUsername.Valid {
			msg.Comments = append(msg.Comments, CommentDB{
				Username: reactionUsername.String,
				Emoji:    emoji.String,
			})
		}

		// Controllo is_read / is_delivered
		if senderUsername.String == username {
			// SEI il mittente → controlla gli altri
			if chatType == "private" {
				var otherUserId int
				err := db.c.QueryRow(`
				SELECT user_id FROM participants
				WHERE conversation_id = ? AND user_id != ?`,
					convId, userId).Scan(&otherUserId)
				if err == nil {
					var delivered, read bool
					_ = db.c.QueryRow(`
					SELECT is_delivered, is_read
					FROM messageReads
					WHERE message_id = ? AND user_id = ?`,
						id, otherUserId).Scan(&delivered, &read)
					msg.IsDelivered = delivered
					msg.IsRead = read
				}
			} else if chatType == "group" {
				var countUndelivered, countUnread int
				_ = db.c.QueryRow(`
				SELECT COUNT(*) FROM messageReads
				WHERE message_id = ? AND user_id != ? AND is_delivered = FALSE`,
					id, userId).Scan(&countUndelivered)
				_ = db.c.QueryRow(`
				SELECT COUNT(*) FROM messageReads
				WHERE message_id = ? AND user_id != ? AND is_read = FALSE`,
					id, userId).Scan(&countUnread)

				msg.IsDelivered = countUndelivered == 0
				msg.IsRead = countUnread == 0
			}
		} else {
			// NON sei il mittente → controlla per TE
			var delivered, read bool
			_ = db.c.QueryRow(`
			SELECT is_delivered, is_read
			FROM messageReads
			WHERE message_id = ? AND user_id = ?`,
				id, userId).Scan(&delivered, &read)

			msg.IsDelivered = delivered
			msg.IsRead = read
		}
	}

	for _, msg := range messageMap {
		messages = append(messages, *msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		if !messages[i].Timestamp.Valid {
			return false
		}
		if !messages[j].Timestamp.Valid {
			return true
		}
		return messages[i].Timestamp.Time.Before(messages[j].Timestamp.Time)
	})

	_ = rows.Close()

	if err := rows.Err(); err != nil {
		return nil, "errore dopo il ciclo dei messaggi: " + err.Error(), ErrInternal
	}

	return messages, "conversazione caricata", nil
}
