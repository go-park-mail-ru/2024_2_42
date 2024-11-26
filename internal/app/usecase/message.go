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
	return muc.mediaRepo.CreateChat(req)
}

func (muc *MessageUsecaseController) GetUserChats(userID uint64) ([]*models.ChatInfo, error) {
	return muc.mediaRepo.GetUserChats(userID)
}
