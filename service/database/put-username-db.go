package database

func (db *appdbimpl) UpdateUsername(uid int, newUsername string) (string, error) {
	query := "UPDATE users SET username = ? WHERE user_id = ?"
	_, err := db.c.Exec(query, newUsername, uid)
	if err != nil {
		if db.IsUniqueConstraintError(err) {
			return "username already taken", ErrConflict
		}
		return err.Error(), ErrInternal
	}
	return "username successfully changed", nil
}
