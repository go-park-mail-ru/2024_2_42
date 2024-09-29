package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"youpin/internal/errors"
	"youpin/internal/models"
	"youpin/internal/models/response"
)

const (
	respSignUpSuccessMesssage = "Registration successful. Please confirm your email"
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

func userIsAlreadySignedUP(u models.User) error {
	for _, user := range registeredUsers {
		if user.Email == u.Email && user.Password == u.Password {
			return errors.ErrUserAlreadyRegistered
		}
	}

	return nil
}

func getUserID(u models.User) uint64 {
	for _, user := range registeredUsers {
		if user.Email == u.Email && user.Password == u.Password {
			return user.UserID
		}
	}

	return 0
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errors.SendErrorResponse(w, errors.ErrorInfo{
			General: err, Internal: errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	user.Sanitize()

	// Incorrect data given
	if err := user.Valid(); err != nil {
		errors.SendErrorResponse(w, errors.ErrorInfo{
			General: err, Internal: errors.ErrUserDataInvalid,
		})
		return
	}

	// User already registered
	if err := userIsAlreadySignedUP(user); err != nil {
		errors.SendErrorResponse(w, errors.ErrorInfo{
			General: err, Internal: errors.ErrUserAlreadyRegistered,
		})
		return
	}

	user.UserID = nextUserID
	nextUserID++

	regUsrMutex.Lock()
	registeredUsers = append(registeredUsers, user)
	regUsrMutex.Unlock()

	SendSignUpResponse(w, response.SignUpResponse{
		UserId: user.UserID, Message: respSignUpSuccessMesssage,
	})

	fmt.Println(registeredUsers)
}

func SendSignUpResponse(w http.ResponseWriter, sr response.SignUpResponse) {
	respJSON, err := json.Marshal(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(respJSON)
}
