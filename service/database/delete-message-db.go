package database

func (db *appdbimpl) DeleteMessageDB(userId int, convId int, messId int) (string, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return "error starting transaction", ErrInternal
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Verifica se l'utente è nella conversazione
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}

	// Elimina il messaggio
	delQuery := `DELETE FROM messages WHERE message_id = ? AND conversation_id = ? AND user_id = ?`
	_, err = tx.Exec(delQuery, messId, convId, userId)
	if err != nil {
		return "error while deleting the message: " + err.Error(), ErrInternal
	}

	// Controlla se esistono ancora messaggi nella conversazione
	checkQuery := `SELECT COUNT(*) FROM messages WHERE conversation_id = ?`
	var count int
	err = tx.QueryRow(checkQuery, convId).Scan(&count)
	if err != nil {
		return "error checking remaining messages: " + err.Error(), ErrInternal
	}

	if count == 0 {
		// Nessun messaggio rimasto, elimina la conversazione
		_, err = tx.Exec(`DELETE FROM participants WHERE conversation_id = ?`, convId)
		if err != nil {
			return "error deleting participants: " + err.Error(), ErrInternal
		}

		_, err = tx.Exec(`DELETE FROM conversations WHERE conversation_id = ?`, convId)
		if err != nil {
			return "error deleting conversation: " + err.Error(), ErrInternal
		}
		// Pulizia extra (opzionale, difensiva)
		_, _ = tx.Exec(`DELETE FROM messageReads WHERE message_id NOT IN (SELECT message_id FROM messages)`)
		_, _ = tx.Exec(`DELETE FROM reactions WHERE message_id NOT IN (SELECT message_id FROM messages)`)
	} else {
		// Se ci sono altri messaggi, aggiorna il lastMessage
		updateQuery := `
			UPDATE conversations
			SET message_id = (
				SELECT message_id FROM messages
				WHERE conversation_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			)
			WHERE conversation_id = ?`
		_, err = tx.Exec(updateQuery, convId, convId)
		if err != nil {
			return "error updating last message: " + err.Error(), ErrInternal
		}
	}

	if err = tx.Commit(); err != nil {
		return "error committing transaction: " + err.Error(), ErrInternal
	}

	return "message deleted correctly and conversation checked", nil
}
