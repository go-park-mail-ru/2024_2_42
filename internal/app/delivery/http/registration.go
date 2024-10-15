package delivery

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	internal_errors "pinset/internal/errors"
	"pinset/internal/models"
	"pinset/internal/models/response"
)

const (
	respSignUpSuccessMesssage = "Registration successful!"
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
	}
)

func userIsAlreadySignedUP(u models.User) error {
	for _, user := range registeredUsers {
		if user.Email == u.Email {
			return internal_errors.ErrUserAlreadyRegistered
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
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	user.Sanitize()

	// Incorrect data given
	if err := user.Valid(); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserDataInvalid,
		})
		return
	}

	// User already registered
	if err := userIsAlreadySignedUP(user); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserAlreadyRegistered,
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
}

func SendSignUpResponse(w http.ResponseWriter, sr response.SignUpResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func TestableUserIsAlreadySignedUP(u models.User) error {
	return userIsAlreadySignedUP(u)
}

func TestableGetUserID(u models.User) uint64 {
	return getUserID(u)
}
