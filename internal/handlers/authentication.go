package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	internal_errors "youpin/internal/errors"
	"youpin/internal/models"
	"youpin/internal/models/request"
	"youpin/internal/models/response"

	"github.com/golang-jwt/jwt"
)

type userSession struct {
	userID uint64
	token  string
}

const tokenExpirationTime = time.Hour * 72

var (
	SECRET = []byte(os.Getenv("JWT_SECRET"))

	sessionsMutex  = &sync.Mutex{}
	activeSessions = map[uint64]userSession{}

	authUsrMutex    = &sync.Mutex{}
	authorizedUsers = []models.User{}
)

func userHasActiveSession(req request.LoginRequest) bool {
	authUsrMutex.Lock()
	defer authUsrMutex.Unlock()

	var id uint64
	for _, user := range authorizedUsers {
		if user.Email == req.Email && user.Password == req.Password {
			id = user.UserID
			break
		}
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	if _, ok := activeSessions[id]; ok {
		return true
	}

	return false
}

// LogIn authenticates a user by checking provided credentials and returns a session token
//
//	@Summary		Authenticate and receive a session token
//	@Description	Allows users to log in by providing valid credentials (username and password). Returns a session token on success.
//	@Tags			Authentication
//	@Param			request	body		object	true	"Login credentials"
//	@Produce		json
//	@Success		200	{object}	session_token	"Successful authentication, session token returned"
//	@Failure		400	{object}	errors.ErrorResponse	"Missing or invalid request body"
//	@Failure		401	{object}	errors.ErrorResponse	"Invalid username or password"
//	@Router			/login [post]
func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var req request.LoginRequest
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}
	user.Email = req.Email
	user.Password = req.Password
	userID := getUserID(user)

	// User is not registered
	if err = userIsAlreadySignedUP(user); err == nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserIsNotRegistered,
		})
	}

	// Does user already have an active session?
	if userHasActiveSession(req) {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserAlreadyAuthorized,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"login":   req.Email,
		"exp":     time.Now().Add(tokenExpirationTime).Unix(),
	})

	signedToken, err := token.SignedString(SECRET)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantSignSessionToken,
		})
		return
	}

	sessionsMutex.Lock()
	activeSessions[userID] = userSession{
		userID: userID, token: signedToken,
	}

	sessionsMutex.Unlock()

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(72 * time.Hour),
	}

	http.SetCookie(w, cookie)

	sendLogInResponse(w, response.LogInResponse{
		SessionCookie: signedToken,
	})
}

func sendLogInResponse(w http.ResponseWriter, sr response.LogInResponse) {
	respJSON, err := json.Marshal(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(respJSON)
}

// LogOut removes the session token for the user
//
//	@Summary		User logout
//	@Description	Removing a valid session token
//	@Tags			Authentication
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	models.MainBoard???	"Successful logout. Cookie went bad"
//	@Failure		403	{object}	errors.ErrorResponse	"Session token is not allowed"
//	@Failure		500	{object}	errors.ErrorResponse	"Bad server response"
//	@Router			/logout [get]
func LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// Is user authorized?
	_, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("Successful logout"))
}

// IsAuthorized checks if a user is authenticated
//
//	@Summary		Get auth status
//	@Description	Returns the user based on a valid session token in the cookie
//	@Tags			Authentication
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	models.User	"Successfully get user details"
//	@Failure		403	{object}	errors.ErrorResponse	"Bad authorization"
//	@Failure		500	{object}	errors.ErrorResponse	"Bad server response"
//	@Router			/is_authorized [get]
func IsAuthorized(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	// Is user authorized?
	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if !token.Valid {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidSessionToken,
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		SendIsAuthResponse(w, response.IsAuthResponse{
			UserID: claims["user_id"].(float64),
		})
	} else {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrBadRequest,
		})
	}
}

func SendIsAuthResponse(w http.ResponseWriter, ar response.IsAuthResponse) {
	respJSON, err := json.Marshal(ar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(respJSON)
}
