package database

func (db *appdbimpl) UpdateImageDB(uid int, newImage string) (string, error) {
	query := "UPDATE users SET url_profile_image = ? WHERE user_id = ?"
	_, err := db.c.Exec(query, newImage, uid)
	if err != nil {
		return "errore nell'impostazione dell'immagine di profilo " + err.Error(), ErrInternal
	}
	return "immagine di profilo cambiata correttamente", nil
}
