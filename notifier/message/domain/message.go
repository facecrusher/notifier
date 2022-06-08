package domain

import "github.com/google/uuid"

const (
	ID_LENGTH = 15
)

type Message struct {
	ID      string `json:"message_id"`
	Message string `json:"message"`
}

func NewMessage(message string) *Message {
	return &Message{
		ID:      uuid.New().String(),
		Message: message,
	}
}
