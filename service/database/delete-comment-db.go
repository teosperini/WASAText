package database

func (db *appdbimpl) DeleteEmojiDB(userId int, convId int, messId int) (string, error) {
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}

	query := `
		DELETE FROM reactions
		WHERE message_id = ? AND user_id = ?
	`
	_, err = db.c.Exec(query, messId, userId)
	if err != nil {
		return "error deleting emoji reaction: " + err.Error(), ErrInternal
	}

	return "", nil
}
