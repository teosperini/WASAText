package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) PostMessageDB(senderId int, convId int, mess MessageToServerDB) (int, string, error) {
	// Inizio della transazione
	tx, err := db.c.Begin()
	if err != nil {
		return 0, "error starting the transaction: " + err.Error(), ErrInternal
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // Se c'è un panic, rollback sicuro
			panic(p)
		} else if err != nil {
			_ = tx.Rollback() // Se c'è un errore, rollback sicuro
		}
	}()

	str, err := db.IsUserInTheConversation(convId, senderId)
	if err != nil {
		return 0, str, err
	}

	fmt.Println("ReplyToMessageId:", mess.ReplyToMessageId)

	// Inserimento del messaggio in base al suo tipo
	var replyId interface{}
	if mess.ReplyToMessageId != nil {
		replyId = *mess.ReplyToMessageId
	} else {
		replyId = nil
	}

	var query string
	var result sql.Result

	switch mess.MessageType {
	case "text":
		if mess.Text == "" {
			err = errors.New("text message is required")
			return 0, "text is required", ErrBadRequest
		}
		query = `
			INSERT INTO messages (conversation_id, user_id, type, content_text, reply_message_id)
			VALUES (?, ?, ?, ?, ?)`
		result, err = tx.Exec(query, convId, senderId, mess.MessageType, mess.Text, replyId)
	case "image":
		if mess.MediaUrl == "" {
			err = errors.New("media url is required")
			return 0, "image is required", ErrBadRequest
		}
		query = `
			INSERT INTO messages (conversation_id, user_id, type, content_image, reply_message_id)
			VALUES (?, ?, ?, ?, ?)`
		result, err = tx.Exec(query, convId, senderId, mess.MessageType, mess.MediaUrl, replyId)
	case "text_image":
		if mess.MediaUrl == "" && mess.Text != "" {
			err = errors.New("text_image is required")
			return 0, "text image is required", ErrBadRequest
		}
		query = `
			INSERT INTO messages (conversation_id, user_id, type, content_text, content_image, reply_message_id)
			VALUES (?, ?, ?, ?, ?, ?)`
		result, err = tx.Exec(query, convId, senderId, mess.MessageType, mess.Text, mess.MediaUrl, replyId)
	default:
		err = errors.New("non-existent message type: " + mess.MessageType)
		return 0, "error: non-existent message type: " + mess.MessageType, ErrBadRequest
	}

	if err != nil {
		return 0, "error sending the message: " + err.Error(), ErrInternal
	}

	// Recupera l'ID del messaggio appena inserito
	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, "error retrieving the message id: " + err.Error(), ErrInternal
	}

	// Aggiorna la conversazione con l'ultimo messaggio inviato
	query = `
		UPDATE conversations
		SET message_id = ?
		WHERE conversation_id = ?`
	_, err = tx.Exec(query, messageId, convId)
	if err != nil {
		return 0, "error updating conversation with last message: " + err.Error(), ErrInternal
	}

	// Inserisce la riga per il sender con delivered e read a true
	query = `
	INSERT INTO messageReads (message_id, user_id, is_delivered, is_read)
	VALUES (?, ?, ?, ?)`
	if _, err := tx.Exec(query, messageId, senderId, true, true); err != nil {
		return 0, "error inserting messageReads for sender: " + err.Error(), ErrInternal
	}

	// Inserisce le righe per tutti gli altri utenti con delivered e read a false
	query = `
	INSERT INTO messageReads (message_id, user_id, is_delivered, is_read)
	SELECT ?, user_id, false, false
	FROM participants
	WHERE conversation_id = ? AND user_id <> ?`

	if _, err := tx.Exec(query, messageId, convId, senderId); err != nil {
		return 0, "error inserting messageReads for recipients: " + err.Error(), ErrInternal
	}

	// Commit della transazione
	err = tx.Commit()
	if err != nil {
		return 0, "error committing the transaction: " + err.Error(), ErrInternal
	}

	return int(messageId), "message sent successfully", nil
}
