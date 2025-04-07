package database

func (db *appdbimpl) UpdateNameGroupDB(userId int, convId int, newName string) (string, error) {
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}
	query := "UPDATE conversations SET group_name = ? WHERE conversation_id = ?"
	_, err = db.c.Exec(query, newName, convId)
	if err != nil {
		return "error while changing the name of the group " + err.Error(), ErrInternal
	}
	return "name of the group successfully changed", nil
}
