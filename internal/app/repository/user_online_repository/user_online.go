package UserOnlineRepository

import (
	"pinset/internal/app/models"
	"sync"
)

type UserOnlineRepositoryController struct {
	mu   *sync.RWMutex
	data map[uint64]*models.ChatUser
}

func NewUserOnlineRepository() *UserOnlineRepositoryController {
	return &UserOnlineRepositoryController{
		mu:   &sync.RWMutex{},
		data: make(map[uint64]*models.ChatUser),
	}
}

func (uoc *UserOnlineRepositoryController) IsOnlineUser(userID uint64) bool {
	uoc.mu.RLock()
	defer uoc.mu.RUnlock()
	_, ok := uoc.data[userID]
	return ok
}

func (uoc *UserOnlineRepositoryController) GetOnlineUser(userID uint64) *models.ChatUser {
	uoc.mu.RLock()
	defer uoc.mu.RUnlock()
	return uoc.data[userID]
}

func (uoc *UserOnlineRepositoryController) AddOnlineUser(user *models.ChatUser) {
	uoc.mu.Lock()
	defer uoc.mu.Unlock()
	uoc.data[user.ID] = user
}

func (uoc *UserOnlineRepositoryController) DeleteOnlineUser(userID uint64) {
	uoc.mu.Lock()
	defer uoc.mu.Unlock()
	delete(uoc.data, userID)
}
