package usecase

import (
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"
)

func NewMessageUsecase(userOnlineRepo UserOnlineRepo, mediaRepo MediaRepository, userRepo UserRepository) delivery.MessageUsecase {
	return &MessageUsecaseController{
		mediaRepo:      mediaRepo,
		userRepo:       userRepo,
		userOnlineRepo: userOnlineRepo,
	}
}

func (muc *MessageUsecaseController) AddOnlineUser(user *models.ChatUser) {
	muc.userOnlineRepo.AddOnlineUser(user)
}

func (muc *MessageUsecaseController) IsOnlineUser(userID uint64) bool {
	return muc.userOnlineRepo.IsOnlineUser(userID)
}

func (muc *MessageUsecaseController) NumUsersOnline() int {
	return muc.userOnlineRepo.NumUsersOnline()
}

func (muc *MessageUsecaseController) GetOnlineUser(userID uint64) *models.ChatUser {
	return muc.userOnlineRepo.GetOnlineUser(userID)
}

func (muc *MessageUsecaseController) DeleteOnlineUser(userID uint64) {
	muc.userOnlineRepo.DeleteOnlineUser(userID)
}

func (muc *MessageUsecaseController) GetChatMessages(chatID uint64) ([]*models.MessageInfo, error) {
	return muc.mediaRepo.GetChatMessages(chatID)
}

func (muc *MessageUsecaseController) AddChatMessage(message *models.Message) (*models.MessageCreateInfo, error) {
	return muc.mediaRepo.CreateMessage(message)
}
func (muc *MessageUsecaseController) GetChatUsers(chatID uint64) ([]uint64, error) {
	return muc.mediaRepo.GetChatUsers(chatID)
}

func (muc *MessageUsecaseController) CreateChat(req *models.ChatCreateRequest) (*models.ChatInfo, error) {
	chatCreateInfo, err := muc.mediaRepo.CreateChat()
	if err != nil {
		return nil, err
	}
	chatID := chatCreateInfo.ID
	err = muc.mediaRepo.AddUserToChat(chatID, req.UserID)
	if err != nil {
		return nil, err
	}
	err = muc.mediaRepo.AddUserToChat(chatID, req.CompanionID)
	if err != nil {
		return nil, err
	}
	companionInfo, err := muc.userRepo.GetUserInfoPublic(req.CompanionID)
	if err != nil {
		return nil, err
	}
	return &models.ChatInfo{ChatID: chatID, Companion: *companionInfo}, nil
}

func (muc *MessageUsecaseController) GetUserChats(userID uint64) ([]*models.ChatInfo, error) {
	chatIDs, err := muc.mediaRepo.GetUserChats(userID)
	if err != nil {
		return nil, err
	}
	chats := make([]*models.ChatInfo, 0)
	for _, chatID := range chatIDs {
		chat := &models.ChatInfo{ChatID: chatID}
		userIDs, err := muc.mediaRepo.GetChatUsers(chatID)
		if err != nil {
			return nil, err
		}
		for _, id := range userIDs {
			if id != userID {
				companion, err := muc.userRepo.GetUserInfoPublic(id)
				if err != nil {
					return nil, err
				}
				chat.Companion = *companion
			}
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
