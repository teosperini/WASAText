package database

func (db *appdbimpl) AddMemberToGroupDB(userId int, convId int, addUsername string) (string, error) {
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}

	// Cerco l'userId dell'utente da aggiungere
	newUserId, err := db.GetUidFromUsername(addUsername)
	if err != nil {
		return "error retrieving the user_id from the username: " + err.Error(), err
	}

	// Aggiungo l'utente
	query := `
		INSERT INTO participants (conversation_id, user_id)
		VALUES (?, ?)`
	_, err = db.c.Exec(query, convId, newUserId)
	if err != nil {
		if db.IsUniqueConstraintError(err) {
			return "user already in the group", ErrConflict
		}
		return err.Error(), ErrInternal
	}
	return "username successfully changed", nil
}
