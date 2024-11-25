package usecase

import (
	"pinset/mailer-service/mailer"
)

type MessageRepository interface {
	GetChatMessages(chatID uint64) ([]*mailer.MessageInfo, error)
	AddChatMessage(message *mailer.Message) (*mailer.MessageInfo, error)
	GetChatUsers(chatID uint64) ([]*mailer.User, error)
	GetUserChats(userID uint64) ([]uint64, error)
	CreateChat() (*mailer.ChatInfo, error)
	AddUserToChat(chatID uint64, userID uint64) error
}

type UserRepository interface {
	GetUserInfoPublic(userID uint64) (*mailer.UserProfileResponse, error)
}

type MessageUsecaseController struct {
	messageRepo MessageRepository
	userRepo    UserRepository
}
