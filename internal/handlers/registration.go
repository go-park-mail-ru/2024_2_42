package handlers

import (
	"net/http"
	"sync"
	"time"

	"youpin/internal/errors"
	"youpin/internal/models"
)

var (
	nextUserID uint64 = 2

	regUsrMutex                   = &sync.Mutex{}
	registeredUsers []models.User = []models.User{
		{
			UserID:       1,
			UserName:     "admin",
			NickName:     "admin",
			Email:        "example@test.com",
			Password:     "12345678",
			BirthTime:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
			Gender:       "table",
			AvatarUrl:    "",
			Followers:    []models.User{},
			Following:    []models.User{},
			Boards:       []models.Board{},
			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
		},
	}
)

func userIsAlreadyRegistered(u models.User) error {
	for _, user := range registeredUsers {
		if user.Email == u.Email {
			return errors.ErrorUserAlreadyRegistered
		}
	}

	return nil
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := models.NewUser(nextUserID, userName, email, password)
	user.Sanitize()

	// Incorrect data given
	if err := user.Valid(); err != nil {
		http.Error(w, err.Error(), errors.ErrorToHttpStatusCode[err])
		return
	}

	// User already registered
	if err := userIsAlreadyRegistered(user); err != nil {
		http.Error(w, err.Error(), errors.ErrorToHttpStatusCode[err])
		return
	}

	regUsrMutex.Lock()
	registeredUsers = append(registeredUsers, user)
	regUsrMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successfull"))
}
