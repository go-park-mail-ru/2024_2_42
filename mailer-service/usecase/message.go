package usecase

import (
	"pinset/mailer-service/delivery"
	"pinset/mailer-service/mailer"
)

func NewMessageUsecase(messageRepo MessageRepository, userRepo UserRepository) delivery.MessageUsecase {
	return &MessageUsecaseController{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (muc *MessageUsecaseController) GetChatMessages(chatID uint64) ([]*mailer.MessageInfo, error) {
	return muc.messageRepo.GetChatMessages(chatID)
}

func (muc *MessageUsecaseController) AddChatMessage(message *mailer.Message) (*mailer.MessageInfo, error) {
	return muc.messageRepo.AddChatMessage(message)
}

func (muc *MessageUsecaseController) GetChatUsers(chatID uint64) ([]*mailer.User, error) {
	return muc.messageRepo.GetChatUsers(chatID)
}

func (muc *MessageUsecaseController) CreateChat(req *mailer.ChatCreateRequest) (*mailer.ChatInfo, error) {
	chatCreateInfo, err := muc.messageRepo.CreateChat()
	if err != nil {
		return nil, err
	}
	chatID := chatCreateInfo.ChatId
	err = muc.messageRepo.AddUserToChat(chatID, req.UserId)
	if err != nil {
		return nil, err
	}
	err = muc.messageRepo.AddUserToChat(chatID, req.CompanionId)
	if err != nil {
		return nil, err
	}
	companionInfo, err := muc.userRepo.GetUserInfoPublic(req.CompanionId)
	if err != nil {
		return nil, err
	}
	return &mailer.ChatInfo{ChatId: chatID, Companion: companionInfo}, nil
}

func (muc *MessageUsecaseController) GetUserChats(userID uint64) ([]*mailer.ChatInfo, error) {
	chatIDs, err := muc.messageRepo.GetUserChats(userID)
	if err != nil {
		return nil, err
	}
	chats := make([]*mailer.ChatInfo, 0)
	for _, chatID := range chatIDs {
		chat := &mailer.ChatInfo{ChatId: chatID}
		users, err := muc.messageRepo.GetChatUsers(chatID)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			if user.Id != userID {
				companion, err := muc.userRepo.GetUserInfoPublic(user.Id)
				if err != nil {
					return nil, err
				}
				chat.Companion = companion
			}
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
