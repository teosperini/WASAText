package database

func (db *appdbimpl) RemoveMemberFromGroupDB(userId int, convId int) (string, error) {

	query := `DELETE from participants WHERE conversation_id = ? AND user_id = ?`
	_, err := db.c.Exec(query, convId, userId)
	if err != nil {
		return "error the user does not participate at the conversation" + err.Error(), ErrForbidden
	}

	return "", nil
}
