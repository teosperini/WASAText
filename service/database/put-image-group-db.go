package database

func (db *appdbimpl) UpdateImageGroupDB(userId int, newImage string, convId int) (string, error) {
	str, err := db.IsUserInTheConversation(convId, userId)
	if err != nil {
		return str, err
	}

	query := "UPDATE conversations SET group_image = ? WHERE conversation_id = ?"
	_, err = db.c.Exec(query, newImage, convId)
	if err != nil {
		return "errore nell'impostazione dell'immagine di profilo del gruppo " + err.Error(), ErrInternal
	}
	return "immagine di profilo del gruppo cambiata correttamente", nil
}
