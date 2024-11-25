package models

import (
	"pinset/internal/app/models/response"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID  uint64    `json:"sender_id"`
	ChatID    uint64    `json:"chat_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type WebSocketResponse struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type MessageInfo struct {
	ID        uint64    `json:"message_id"`
	SenderID  uint64    `json:"sender_id"`
	ChatID    uint64    `json:"chat_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorInfo struct {
	Code    int `json:"error_code"`
	Message int `json:"error_message"`
}

type MessageUpdate struct {
	ID      uint64 `json:"message_id"`
	Content string `json:"content"`
}

type MessageCreateInfo struct {
	ID        uint64    `json:"message_id"`
	SenderID  uint64    `json:"sender_id"`
	ChatID    uint64    `json:"chat_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatCreateInfo struct {
	ID uint64 `json:"chat_id"`
}

type ChatJoiner struct {
	ChatID uint64 `json:"chat_id"`
}

type ChatUser struct {
	ID         uint64
	Connection *websocket.Conn
}

type ChatInfo struct {
	ChatID    uint64                       `json:"chat_id"`
	Companion response.UserProfileResponse `json:"companion"`
}

type ChatCreateRequest struct {
	UserID      uint64
	CompanionID uint64
}
