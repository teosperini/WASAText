package database

import (
	"errors"
)

// Recupera o crea un utente dato l'username
func (db *appdbimpl) DoLoginDB(username string, image string) (int, string, error) {
	// Controlla se l'utente esiste
	userId, err := db.GetUidFromUsername(username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// L'utente non esiste, quindi lo crea
			result, err := db.c.Exec("INSERT INTO users (username, url_profile_image) VALUES (?, ?)", username, image)
			if err != nil {
				return 0, "error inserting user into the users table: " + err.Error(), ErrInternal
			}
			// Ottiene l'ID del nuovo utente
			insertedId, err := result.LastInsertId()
			if err != nil {
				return 0, "error retrieving the id of the last inserted user: " + err.Error(), ErrInternal
			}

			userId = int(insertedId)

			return userId, "User " + username + "created and userId returned", nil
		} else {
			return 0, "error is in the GetUidFromUsername function: " + err.Error(), ErrInternal
		}
	}

	return userId, "User " + username + " found and userId returned", nil
}
