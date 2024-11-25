package delivery

import (
	"context"
	"pinset/mailer-service/mailer"
)

type MessageUsecase interface {
	GetChatMessages(chatID uint64) ([]*mailer.MessageInfo, error)
	AddChatMessage(message *mailer.Message) (*mailer.MessageInfo, error)
	GetChatUsers(chatID uint64) ([]*mailer.User, error)
	GetUserChats(userID uint64) ([]*mailer.ChatInfo, error)
	CreateChat(req *mailer.ChatCreateRequest) (*mailer.ChatInfo, error)
}

type MessageDelieveryController struct {
	mailer.UnimplementedChatServiceServer
	Usecase MessageUsecase
}

func NewMessageDeliveryController(usecase MessageUsecase) *MessageDelieveryController {
	return &MessageDelieveryController{
		Usecase: usecase,
	}
}

func (mc *MessageDelieveryController) AddChatMessage(ctx context.Context, msg *mailer.Message) (*mailer.MessageInfo, error) {
	createInfo, err := mc.Usecase.AddChatMessage(msg)
	if err != nil {
		return nil, err
	}
	return createInfo, nil
}
func (mc *MessageDelieveryController) CreateChat(ctx context.Context, req *mailer.ChatCreateRequest) (*mailer.ChatInfo, error) {
	chatInfo, err := mc.Usecase.CreateChat(req)
	if err != nil {
		return nil, err
	}
	return chatInfo, nil
}
func (mc *MessageDelieveryController) GetAllChatMessages(ctx context.Context, req *mailer.ChatRequest) (*mailer.MessageList, error) {
	messages, err := mc.Usecase.GetChatMessages(req.ChatId)
	if err != nil {
		return nil, err
	}
	return &mailer.MessageList{Messages: messages}, nil
}
func (mc *MessageDelieveryController) GetChatUsers(ctx context.Context, req *mailer.ChatRequest) (*mailer.UserList, error) {
	users, err := mc.Usecase.GetChatUsers(req.ChatId)
	if err != nil {
		return nil, err
	}
	return &mailer.UserList{Users: users}, nil
}
func (mc *MessageDelieveryController) GetUserChats(ctx context.Context, user *mailer.User) (*mailer.ChatList, error) {
	chats, err := mc.Usecase.GetUserChats(user.Id)
	if err != nil {
		return nil, err
	}
	return &mailer.ChatList{Chats: chats}, nil
}
