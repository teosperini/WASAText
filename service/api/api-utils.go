package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/teosperini/WASAText/service/database"
)

const (
	ErrWriteResponse = "failed to write response: "
	PrefixUser       = "User "
	InvalidId        = "invalid id: "
)

func checkLength(username string) bool {
	if len(username) < 3 || len(username) > 16 {
		return true
	}
	return false
}

// This function converts the conversation to create to the DB type
func convCreateToDB(conv ConversationCreateAPI) database.ConversationCreateDB {
	return database.ConversationCreateDB{
		Members:  conv.Members,
		ChatType: conv.ChatType,
		InitialMessage: database.MessageToServerDB{
			MessageType:          conv.InitialMessage.MessageType,
			Text:                 conv.InitialMessage.Text,
			MediaUrl:             conv.InitialMessage.MediaUrl,
			ForwardFromMessageId: conv.InitialMessage.ForwardFromMessageId,
		},
		GroupImageUrl: conv.GroupImageUrl,
		GroupName:     conv.GroupName,
	}
}

func messCreateToDB(mess MessageToServerAPI) database.MessageToServerDB {
	return database.MessageToServerDB{
		MessageType:      mess.MessageType,
		Text:             mess.Text,
		MediaUrl:         mess.MediaUrl,
		ReplyToMessageId: mess.ReplyToMessageId,
	}
}

// Extract the token in the Authorization header and asks the database if it exists
func (rt *_router) extractId(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, database.ErrUnauthorized
	}

	ID := strings.TrimPrefix(authHeader, "Bearer ")
	parsedID, err := strconv.Atoi(ID)
	if err != nil {
		return 0, database.ErrUnauthorized
	}

	err = rt.db.ValidateID(parsedID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, database.ErrNotFound
		} else {
			return 0, database.ErrInternal
		}
	}

	return parsedID, nil
}

func checkErrors(err error) (string, int) {
	if errors.Is(err, database.ErrBadRequest) {
		return "Bad Request - Invalid input", http.StatusBadRequest
	} else if errors.Is(err, database.ErrNotFound) {
		return "Not Found - Please try again", http.StatusNotFound
	} else if errors.Is(err, database.ErrUnauthorized) {
		return "Unauthorized - Please authenticate", http.StatusUnauthorized
	} else if errors.Is(err, database.ErrConflict) {
		return "Conflict - Please try again", http.StatusConflict
	} else if errors.Is(err, database.ErrForbidden) {
		return "Forbidden - Access denied", http.StatusForbidden
	}
	return "Internal Server Error - Please try again later", http.StatusInternalServerError
}

// Converte una conversazione del database in una preview per l'API
func ConvertToConversationPreview(conv database.ConversationDB) ConversationPreview {
	return ConversationPreview{
		ConversationID: conv.ConvID,
		ChatName:       conv.ChatName,
		ChatImageUrl:   conv.ChatImageUrl,
		ChatType:       conv.ChatType,
		UnreadMessages: conv.UnreadMessages,
		Members:        conv.Members,
		LastMessage:    ConvertToMessageToClient(conv.LastMessage),
	}
}

// Converte una conversazione completa in ConversationMessages per l'API

func ConvertToConversationMessages(convID int, messages []database.MessageDB) ConversationMessages {
	var messageList []MessageToClient
	for _, msg := range messages {
		messageList = append(messageList, ConvertToMessageToClient(msg))
	}
	return ConversationMessages{
		ConversationID: convID,
		Messages:       messageList,
	}
}

// Converte un messaggio dal DB al formato API

func ConvertToMessageToClient(msg database.MessageDB) MessageToClient {
	var commentList []Comment
	if len(msg.Comments) > 0 {
		for _, c := range msg.Comments {
			commentList = append(commentList, Comment{
				Username: c.Username,
				Emoji:    c.Emoji,
			})
		}
	}
	if commentList == nil {
		commentList = []Comment{}
	}
	return MessageToClient{
		MessageID:      msg.MessageID,
		SenderUsername: msg.SenderUsername,
		Timestamp:      msg.Timestamp.Time.Format(time.RFC3339),
		MessageType:    msg.MessageType,
		Text:           msg.Text,
		ImageUrl:       msg.ImageUrl,
		IsDelivered:    msg.IsDelivered,
		IsRead:         msg.IsRead,
		IsForwarded:    msg.IsForwarded,
		IsAnswering:    msg.IsAnswering,
		Comments:       commentList,
	}
}

func toEmojiDB(emoji EmojiAPI) database.EmojiDB {
	return database.EmojiDB{
		Emoji: emoji.Emoji,
	}
}
