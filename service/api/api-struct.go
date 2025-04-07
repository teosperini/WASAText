package api

type UserAPI struct {
	Username string `json:"username"`
	Avatar   string `json:"profileImageUrl"`
}

type Username struct {
	Username string `json:"username"`
}

type GroupName struct {
	GroupName string `json:"groupName"`
}

type GroupImage struct {
	GroupImage string `json:"chatImageUrl"`
}

type Image struct {
	ProfileImageUrl string `json:"profileImageUrl"`
}

// Create conversation

type ConversationCreateAPI struct {
	Members        []string           `json:"members"`
	ChatType       string             `json:"chatType"`
	InitialMessage MessageToServerAPI `json:"initialMessage"`
	GroupImageUrl  string             `json:"groupImageUrl,omitempty"`
	GroupName      string             `json:"groupName,omitempty"`
}

// The first message of the created conversation

type MessageToServerAPI struct {
	MessageType          string `json:"messageType"`
	Text                 string `json:"text,omitempty"`
	MediaUrl             string `json:"mediaUrl,omitempty"`
	ForwardFromMessageId *int   `json:"forwardFromMessageId,omitempty"`
	ReplyToMessageId     *int   `json:"replyToMessageId,omitempty"`
}

// Return conversation
// Struttura per la preview della conversazione

type ConversationPreview struct {
	ConversationID int             `json:"conversationId"`
	ChatName       string          `json:"chatName"`
	ChatImageUrl   string          `json:"chatImageUrl"`
	ChatType       string          `json:"chatType"`
	UnreadMessages int             `json:"unreadMessages"`
	Members        []string        `json:"members"`
	LastMessage    MessageToClient `json:"lastMessage"`
}

// Struttura per i messaggi di una conversazione

type ConversationMessages struct {
	ConversationID int               `json:"conversationId"`
	Messages       []MessageToClient `json:"messages"`
}

// Struttura di ritorno per un messaggio da inviare al client

type MessageToClient struct {
	MessageID      int       `json:"messageId"`
	SenderUsername string    `json:"senderUsername"`
	Timestamp      string    `json:"timestamp"`
	MessageType    string    `json:"messageType"`
	IsDelivered    bool      `json:"isDelivered"`
	IsRead         bool      `json:"isRead"`
	Text           string    `json:"text,omitempty"`
	ImageUrl       string    `json:"mediaUrl,omitempty"`
	IsForwarded    bool      `json:"isForwarded"`
	IsAnswering    int       `json:"isAnswering,omitempty"`
	Comments       []Comment `json:"comments"`
}

// Struttura di ritorno per un commento al client

type Comment struct {
	Username string `json:"username"`
	Emoji    string `json:"emoji"`
}

type EmojiAPI struct {
	Emoji string `json:"emoji"`
}
