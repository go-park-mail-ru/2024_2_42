package repository

import (
	"pinset/internal/app/models"
	"pinset/internal/app/session"
	"pinset/internal/app/usecase"
	internal_errors "pinset/internal/errors"
	"sync"
	"time"
)

var (
	nextUserID uint64 = 2
)

func NewUserRepository() usecase.UserRepository {
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
		mu: &sync.RWMutex{},
		sm: session.NewSessionManager(),
	}
}

func (urc *UserRepositoryController) Insert(user *models.User) error {
	user.UserID = nextUserID
	nextUserID++

	urc.mu.Lock()
	defer urc.mu.Unlock()

	// User is already in database
	if _, ok := urc.db[user.Email]; ok {
		return internal_errors.ErrUserAlreadyRegistered
	}

	urc.db[user.Email] = user

	return nil
}

func (urc *UserRepositoryController) UserAlreadySignedUp(user models.User) bool {
	urc.mu.Lock()
	defer urc.mu.Unlock()

	u, ok := urc.db[user.Email]

	return ok && u.Password == user.Password
}

func (urc *UserRepositoryController) GetUserId(user models.User) uint64 {
	urc.mu.Lock()
	defer urc.mu.Unlock()

	if dbUser, ok := urc.db[user.Email]; !ok {
		return 0
	} else {
		return dbUser.UserID
	}
}

func (urc *UserRepositoryController) UserHasActiveSession(token string) bool {
	return urc.sm.Exists(token)
}

func (urc *UserRepositoryController) Session() *session.SessionsManager {
	return urc.sm
}
