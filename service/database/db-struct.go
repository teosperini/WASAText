package database

import "database/sql"

type UserDB struct {
	UserID   int
	Username string
	Image    string
}

type ConversationCreateDB struct {
	Members        []string
	ChatType       string
	InitialMessage MessageToServerDB
	GroupImageUrl  string
	GroupName      string
}

type MessageToServerDB struct {
	MessageType          string
	Text                 string
	MediaUrl             string
	ForwardFromMessageId *int
	ReplyToMessageId     *int
}

// Struttura del database per una conversazione

type ConversationDB struct {
	ConvID         int
	ChatType       string
	ChatName       string
	ChatImageUrl   string
	UnreadMessages int
	Members        []string // Array con i nomi degli utenti
	LastMessage    MessageDB
}

// Struttura del database per un messaggio

type MessageDB struct {
	MessageID      int
	SenderUsername string
	Timestamp      sql.NullTime
	MessageType    string
	Text           string
	ImageUrl       string
	IsDelivered    bool
	IsRead         bool
	IsForwarded    bool
	IsAnswering    int
	Comments       []CommentDB
}

// Struttura per un commento su un messaggio

type CommentDB struct {
	Username string
	Emoji    string
}

type EmojiDB struct {
	Emoji string
}
