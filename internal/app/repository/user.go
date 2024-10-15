package repository

import (
	"pinset/internal/models"
	"time"
)

func NewUserRepository() UserRepository {
	return &UserRepositoryController{
		db: map[string]*models.User{
			"example@test.com": {
				UserID:       1,
				UserName:     "admin",
				NickName:     "admin",
				Email:        "example@test.com",
				Password:     "12345678Q",
				BirthTime:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
				Gender:       "table",
				AvatarUrl:    "",
				Followers:    []models.User{},
				Following:    []models.User{},
				Boards:       []models.Board{},
				CreationTime: time.Now(),
				UpdateTime:   time.Now(),
			},
		},
	}
}

func (urc *UserRepositoryController) UserAlreadySignedUp(u models.User) bool {
	urc.mu.Lock()
	defer urc.mu.Unlock()

	_, ok := urc.db[u.Email]

	return ok
}

func (urc *UserRepositoryController) GetUserId(u models.User) uint64 {
	urc.mu.Lock()
	defer urc.mu.Unlock()

	if dbUser, ok := urc.db[u.Email]; !ok {
		return 0
	} else {
		return dbUser.UserID
	}
}

func (urc *UserRepositoryController) UserHasActiveSession(token string) bool {
	return urc.sm.Exists(token)
}
