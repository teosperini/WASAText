package database

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) PostConversationDB(creatorId int, conv ConversationCreateDB) (int, []MessageDB, string, error) {
	if conv.InitialMessage.ForwardFromMessageId != nil {
		originalMsg, err := db.GetMessageById(*conv.InitialMessage.ForwardFromMessageId)
		if err != nil {
			return 0, nil, "original message not found: " + err.Error(), ErrNotFound
		}

		conv.InitialMessage.Text = originalMsg.Text
		conv.InitialMessage.MediaUrl = originalMsg.ImageUrl
		conv.InitialMessage.MessageType = originalMsg.MessageType
	}

	if conv.InitialMessage.Text == "" && conv.InitialMessage.MediaUrl == "" {
		return 0, nil, "error: initial message must have valid content", ErrBadRequest
	}

	if conv.ChatType == "private" && len(conv.Members) != 1 {
		return 0, nil, "error, can't create a chat with more than two members", ErrBadRequest
	}

	if conv.ChatType == "group" {
		if conv.GroupName == "" {
			return 0, nil, "error, can't create a group chat without a name", ErrBadRequest
		}
		if len(conv.Members) < 1 {
			return 0, nil, "error, can't create a group chat with less than one members", ErrBadRequest
		}
	}

	// Inizio della transazione
	tx, err := db.c.Begin()
	if err != nil {
		return 0, nil, "error starting the transaction: " + err.Error(), ErrInternal
	}

	// Rollback automatico in caso di errore
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Error("rollback failed:", rollbackErr)
			}
		}
	}()

	var result sql.Result

	if conv.ChatType == "private" {
		query := `
			SELECT c.conversation_id
			FROM conversations c
			INNER JOIN participants cp1 ON c.conversation_id = cp1.conversation_id
			INNER JOIN participants cp2 ON c.conversation_id = cp2.conversation_id
			WHERE c.type = 'private'
			AND cp1.user_id = ?
			AND cp2.user_id = ?;`

		var conversationId int
		uid, err := db.GetUidFromUsername(conv.Members[0])
		if err != nil {
			return 0, nil, "error retrieving the id for the other user", err
		}

		// Controllo se esiste già una conversazione privata tra i due utenti
		err = db.c.QueryRow(query, creatorId, uid).Scan(&conversationId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Nessuna chat trovata, quindi possiamo crearne una nuova
			} else {
				// Se c'è un altro errore, restituiscilo
				return 0, nil, "database error while checking for existing chat: " + err.Error(), ErrInternal
			}
		} else {
			if creatorId == uid {
				return 0, nil, "a private conversation between the same user and itself can't be created", ErrBadRequest
			}
			// Se troviamo una conversazione esistente, restituiamo un errore
			return 0, nil, "a private conversation between the users already exists", ErrBadRequest
		}

		query = `
			INSERT INTO conversations (type)
			VALUES (?)`
		result, err = tx.Exec(query, conv.ChatType)
		if err != nil {
			return 0, nil, "error during the creation of a private chat: " + err.Error(), ErrBadRequest
		}
	} else {
		logrus.Println(conv.GroupImageUrl)
		if conv.GroupImageUrl == "" {
			conv.GroupImageUrl = "https://cdn-icons-png.flaticon.com/128/11989/11989096.png"
		}

		query := `
			INSERT INTO conversations (type, group_name, group_image)
			VALUES (?, ?, ?)`
		result, err = tx.Exec(query, conv.ChatType, conv.GroupName, conv.GroupImageUrl)
		if err != nil {
			return 0, nil, "error inserting the conversation into the database: " + err.Error(), ErrBadRequest
		}
	}

	// Ottieni l'ID della conversazione appena creata
	convId, err := result.LastInsertId()
	if err != nil {
		return 0, nil, "error retrieving the id of the conversation: " + err.Error(), ErrInternal
	}

	var mt = conv.InitialMessage.MessageType
	var msgResult sql.Result

	// Inserimento del messaggio iniziale e recupero del message id
	// Inserimento del messaggio iniziale e recupero del message id
	if mt == "text" {
		query := `
			INSERT INTO messages (conversation_id, user_id, type, content_text, is_forwarded)
			VALUES (?, ?, ?, ?, ?)`
		msgResult, err = tx.Exec(query, convId, creatorId, conv.InitialMessage.MessageType, conv.InitialMessage.Text, conv.InitialMessage.ForwardFromMessageId != nil)
	} else if mt == "image" {
		query := `
			INSERT INTO messages (conversation_id, user_id, type, content_image, is_forwarded)
			VALUES (?, ?, ?, ?, ?)`
		msgResult, err = tx.Exec(query, convId, creatorId, conv.InitialMessage.MessageType, conv.InitialMessage.MediaUrl, conv.InitialMessage.ForwardFromMessageId != nil)
	} else if mt == "image_text" {
		query := `
			INSERT INTO messages (conversation_id, user_id, type, content_text, content_image, is_forwarded)
			VALUES (?, ?, ?, ?, ?, ?)`
		msgResult, err = tx.Exec(query, convId, creatorId, conv.InitialMessage.MessageType, conv.InitialMessage.Text, conv.InitialMessage.MediaUrl, conv.InitialMessage.ForwardFromMessageId != nil)
	} else {
		return 0, nil, "error: non-existent message type: " + mt, ErrBadRequest
	}

	if err != nil {
		return 0, nil, "error creating the message: " + err.Error(), ErrInternal
	}

	// Recupera il message id appena inserito
	messageId, err := msgResult.LastInsertId()
	if err != nil {
		return 0, nil, "error retrieving the message id: " + err.Error(), ErrInternal
	}

	// Aggiorna il campo `message_id` nella tabella `conversations`
	query := `
		UPDATE conversations
		SET message_id = (SELECT MAX(message_id) FROM messages WHERE conversation_id = ?)
		WHERE conversation_id = ?`
	_, err = tx.Exec(query, convId, convId)
	if err != nil {
		return 0, nil, "error attaching the message to the conversation: " + err.Error(), ErrInternal
	}

	// Inserisci i partecipanti nella tabella `participants`
	for _, user := range conv.Members {
		query := `
			INSERT INTO participants (conversation_id, user_id)
			VALUES (?, ?)`
		var id int
		id, err = db.GetUidFromUsername(user)
		if err != nil {
			return 0, nil, "error retrieving the user_id from the username: " + err.Error(), err
		}
		_, err = tx.Exec(query, convId, id)
		if err != nil {
			return 0, nil, "error adding \"" + user + "\" to the chat: " + err.Error(), ErrInternal
		}
	}

	// Aggiungi il creatore alla lista dei partecipanti
	query = `
		INSERT INTO participants (conversation_id, user_id)
		VALUES (?, ?)`
	_, err = tx.Exec(query, convId, creatorId)
	if err != nil {
		return 0, nil, "error adding the creator to the chat: " + err.Error(), ErrBadRequest
	}

	// Ora, inseriamo le righe in messageReads per ogni partecipante
	// Per il creatore: il messaggio è già stato ricevuto e letto
	insertMR := `
		INSERT INTO messageReads (message_id, user_id, is_delivered, is_read)
		VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(insertMR, messageId, creatorId, true, true)
	if err != nil {
		return 0, nil, "error inserting messageReads for sender: " + err.Error(), ErrInternal
	}

	// Per gli altri partecipanti: il messaggio non è ancora stato consegnato/letto
	for _, user := range conv.Members {
		var uid int
		uid, err = db.GetUidFromUsername(user)
		if err != nil {
			return 0, nil, "error retrieving the user_id from the username: " + err.Error(), err
		}
		_, err = tx.Exec(insertMR, messageId, uid, false, false)
		if err != nil {
			return 0, nil, "error inserting messageReads for recipient (" + user + "): " + err.Error(), ErrInternal
		}
	}

	// Commit della transazione
	err = tx.Commit()
	if err != nil {
		return 0, nil, "error committing the transaction: " + err.Error(), ErrInternal
	}

	// Verifica che il messageId sia valido prima di eseguire la query
	if messageId == 0 {
		return 0, nil, "error: messageId is zero, message insertion may have failed", ErrInternal
	}

	// time.Sleep(100 * time.Millisecond)

	// Recuperiamo il messaggio appena inserito
	query = `
		SELECT m.message_id, u.username, m.timestamp, m.type, m.content_text, m.content_image, m.is_forwarded
		FROM messages m
		INNER JOIN users u ON m.user_id = u.user_id
		WHERE m.message_id = ?`

	row := db.c.QueryRow(query, messageId)

	var msg MessageDB
	var text sql.NullString
	var imageUrl sql.NullString
	var isForwarded bool

	if err = row.Scan(&msg.MessageID, &msg.SenderUsername, &msg.Timestamp, &msg.MessageType, &text, &imageUrl, &isForwarded); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil, "error: the inserted message was not found in the database", ErrInternal
		}
		return 0, nil, "error retrieving the inserted message: " + err.Error(), ErrInternal
	}

	// Gestiamo i valori NULL
	msg.Text = text.String
	msg.ImageUrl = imageUrl.String
	msg.IsDelivered = false
	msg.IsRead = false
	msg.IsForwarded = isForwarded

	// Creiamo un array con il messaggio appena inserito, per renderlo conforme all'API
	messages := []MessageDB{msg}

	return int(convId), messages, ":D", nil
}
