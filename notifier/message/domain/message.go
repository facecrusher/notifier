package domain

import (
	"github.com/google/uuid"
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
