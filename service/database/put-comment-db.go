package database

func (db *appdbimpl) PutEmojiDB(userId int, convId int, messId int, emoji EmojiDB) (string, error) {
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}

	// Inizio della transazione
	tx, err := db.c.Begin()
	if err != nil {
		return "error starting the transaction: " + err.Error(), ErrInternal
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // Se c'è un panic, rollback sicuro
			panic(p)
		} else if err != nil {
			_ = tx.Rollback() // Se c'è un errore, rollback sicuro
		}
	}()

	query := `INSERT INTO reactions (message_id, user_id, unicode)
		VALUES (?, ?, ?)
		ON CONFLICT(message_id, user_id) DO UPDATE SET unicode = excluded.unicode`
	_, err = tx.Exec(query, messId, userId, emoji.Emoji)
	if err != nil {
		return "error uploading the comment: " + err.Error(), ErrInternal
	}

	err = tx.Commit()
	if err != nil {
		return "error committing the transaction: " + err.Error(), ErrInternal
	}

	return "", nil
}
