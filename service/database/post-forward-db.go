package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) ForwardMessageDB(userId int, newConvId int, messId int) (int, string, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, "error starting the transaction: " + err.Error(), ErrInternal
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	var count int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM participants WHERE conversation_id = ? AND user_id = ?`,
		newConvId, userId,
	).Scan(&count)
	if err != nil {
		return 0, "error checking conversation membership: " + err.Error(), ErrInternal
	}
	if count == 0 {
		err = errors.New("user is not a participant in the conversation")
		return 0, "user is not a participant in the conversation", ErrForbidden
	}

	// Recupera i dati del messaggio originale
	var msg MessageDB
	var convId int
	query := `
		SELECT conversation_id, type, content_text, content_image
		FROM messages
		WHERE message_id = ?`

	err = tx.QueryRow(query, messId).Scan(&convId, &msg.MessageType, &msg.Text, &msg.ImageUrl)

	err = tx.QueryRow(
		`SELECT COUNT(*) FROM participants WHERE conversation_id = ? AND user_id = ?`,
		convId, userId,
	).Scan(&count)
	if err != nil {
		return 0, "error checking conversation membership: " + err.Error(), ErrInternal
	}
	if count == 0 {
		err = errors.New("user is not a participant in the conversation")
		return 0, "user is not a participant in the conversation", ErrForbidden
	}
	// Controlla se l'utente appartiene alla conversazione da dove prende il messaggio

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "messaggio da inoltrare non trovato", ErrNotFound
		}
		return 0, "errore nel recupero del messaggio " + err.Error(), ErrInternal
	}

	// Inserisci il nuovo messaggio (inoltrato)
	var result sql.Result
	if msg.MessageType == "text" {
		result, err = tx.Exec(`
			INSERT INTO messages (conversation_id, user_id, type, content_text, is_forwarded)
			VALUES (?, ?, ?, ?, ?)`,
			newConvId, userId, msg.MessageType, msg.Text, true)
	} else if msg.MessageType == "image" {
		result, err = tx.Exec(`
			INSERT INTO messages (conversation_id, user_id, type, content_image, is_forwarded)
			VALUES (?, ?, ?, ?, ?)`,
			newConvId, userId, msg.MessageType, msg.ImageUrl, true)
	} else if msg.MessageType == "text_image" {
		result, err = tx.Exec(`
			INSERT INTO messages (conversation_id, user_id, type, content_text, content_image, is_forwarded)
			VALUES (?, ?, ?, ?, ?, ?)`,
			newConvId, userId, msg.MessageType, msg.Text, msg.ImageUrl, true)
	} else {
		return 0, "message type not supported for forwarding: " + msg.MessageType, ErrBadRequest
	}

	if err != nil {
		return 0, "errore nell'inserimento del messaggio inoltrato: " + err.Error(), ErrInternal
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, "error retrieving the message id: " + err.Error(), ErrInternal
	}

	// Mittente (utente che inoltra)
	query = `
		INSERT INTO messageReads (message_id, user_id, is_delivered, is_read)
		VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(query, messageId, userId, true, true)
	if err != nil {
		return 0, "errore nell'inserimento di messageReads per il mittente: " + err.Error(), ErrInternal
	}

	// Partecipanti (escluso mittente)
	query = `
		INSERT INTO messageReads (message_id, user_id, is_delivered, is_read)
		SELECT ?, user_id, false, false
		FROM participants
		WHERE conversation_id = ? AND user_id != ?`
	_, err = tx.Exec(query, messageId, newConvId, userId)
	if err != nil {
		return 0, "errore nell'inserimento di messageReads per i destinatari: " + err.Error(), ErrInternal
	}

	query = `
		UPDATE conversations
		SET message_id = ?
		WHERE conversation_id = ?`
	_, err = tx.Exec(query, messageId, newConvId)
	if err != nil {
		return 0, "errore nell'aggiornamento del last message nella conversazione: " + err.Error(), ErrInternal
	}

	err = tx.Commit()
	if err != nil {
		return 0, "errore nel commit della transazione: " + err.Error(), ErrInternal
	}

	return int(messageId), ":D", nil
}
