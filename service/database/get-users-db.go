package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUsersDB(username string) ([]UserDB, string, error) {
	query := "SELECT user_id, username, url_profile_image FROM users WHERE username LIKE ?"
	rows, err := db.c.Query(query, "%"+username+"%")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "no user found with the given name", ErrNotFound
		}
		return nil, "internal error while executing the user search query", ErrInternal
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var userList []UserDB

	for rows.Next() {
		var user UserDB
		if err := rows.Scan(&user.UserID, &user.Username, &user.Image); err != nil {
			return nil, "internal error while scanning user data", ErrInternal
		}
		userList = append(userList, user)
	}
	if err := rows.Err(); err != nil {
		return nil, "errore durante la scan delle righe " + err.Error(), ErrInternal
	}

	if err := rows.Err(); err != nil {
		return nil, "internal error while iterating over the user list", ErrInternal
	}

	return userList, "", nil
}
