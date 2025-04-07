/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	DoLoginDB(username string, image string) (int, string, error)
	Ping() error
	ValidateID(ID int) error
	MarkMessagesAsDelivered(userId int, conversations []ConversationDB) error
	GetUrlFromUid(id int) (string, error)
	GetUsersDB(username string) ([]UserDB, string, error)
	UpdateUsername(uid int, newUsername string) (string, error)
	UpdateImageDB(uid int, newImage string) (string, error)
	IsUniqueConstraintError(err error) bool
	PostConversationDB(creatorId int, conv ConversationCreateDB) (int, []MessageDB, string, error)
	GetConversations(userId int) ([]ConversationDB, string, error)
	PostMessageDB(senderId int, conversationId int, mess MessageToServerDB) (int, string, error)
	GetConversationDB(userId int, convId int) ([]MessageDB, string, error)
	AddMemberToGroupDB(userId int, convId int, addUsername string) (string, error)
	RemoveMemberFromGroupDB(userId int, convId int) (string, error)
	UpdateImageGroupDB(userId int, newImage string, convId int) (string, error)
	UpdateNameGroupDB(userId int, convId int, newName string) (string, error)
	DeleteMessageDB(userId int, convId int, messId int) (string, error)
	ForwardMessageDB(userId int, newConvId int, messId int) (int, string, error)
	PutEmojiDB(userId int, convId int, messId int, emoji EmojiDB) (string, error)
	DeleteEmojiDB(userId int, convId int, messId int) (string, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	_, errPragma := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPragma != nil {
		return nil, errPragma
	}
	// Check if the 'users' table exists
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// Create the 'users' table
		users := `
			CREATE TABLE users (
				user_id INTEGER NOT NULL PRIMARY KEY,       -- Identificativo unico
				username TEXT UNIQUE NOT NULL,         -- Nome utente unico
				url_profile_image TEXT                 -- URL dell'immagine di profilo
			);`

		_, err = db.Exec(users)
		if err != nil {
			return nil, fmt.Errorf("error creating users tabel: %w", err)
		}

		conversations :=
			`CREATE TABLE conversations (
				conversation_id INTEGER PRIMARY KEY AUTOINCREMENT,
				type TEXT CHECK(type IN ('private', 'group')) NOT NULL,
				message_id INTEGER,
				group_name TEXT,
				group_image TEXT,
    			FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE SET NULL,
				CHECK((type = 'private' AND group_name IS NULL AND group_image IS NULL) OR
					  (type = 'group' AND group_name IS NOT NULL))
			);`

		_, err = db.Exec(conversations)
		if err != nil {
			return nil, fmt.Errorf("error creating conversations table: %w", err)
		}

		participants := `CREATE TABLE participants (
			    conversation_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				PRIMARY KEY (conversation_id, user_id),
				FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id),
				FOREIGN KEY (user_id) REFERENCES users(user_id)
			);`

		_, err = db.Exec(participants)
		if err != nil {
			return nil, fmt.Errorf("error creating participants table: %w", err)
		}

		messages := `CREATE TABLE messages (
				message_id INTEGER PRIMARY KEY AUTOINCREMENT,
				conversation_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				type TEXT NOT NULL CHECK (type IN ('text', 'image', 'text_image')),
				timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
				content_text TEXT,
    			content_image TEXT,
    			reply_message_id INTEGER REFERENCES messages(message_id) ON DELETE SET NULL,
    			is_forwarded BOOLEAN NOT NULL DEFAULT FALSE,
				FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id) ON DELETE CASCADE,
				FOREIGN KEY (user_id) REFERENCES users(user_id)
			);`

		_, err = db.Exec(messages)
		if err != nil {
			return nil, fmt.Errorf("error creating messages table: %w", err)
		}

		messageReads := `CREATE TABLE messageReads (
				message_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				is_delivered BOOLEAN NOT NULL DEFAULT FALSE,
				is_read BOOLEAN NOT NULL DEFAULT FALSE,
				PRIMARY KEY (message_id, user_id),
				FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE,
				FOREIGN KEY (user_id) REFERENCES users(user_id)
			);`

		_, err = db.Exec(messageReads)
		if err != nil {
			return nil, fmt.Errorf("error creating message reads table: %w", err)
		}

		reactions := `CREATE TABLE reactions (
				message_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				unicode TEXT NOT NULL,
				PRIMARY KEY (message_id, user_id),
				FOREIGN KEY (message_id) REFERENCES Messages(message_id) ON DELETE CASCADE,
				FOREIGN KEY (user_id) REFERENCES Users(user_id)
			);`

		_, err = db.Exec(reactions)
		if err != nil {
			return nil, fmt.Errorf("error creating message reads table: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
