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
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
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

	if !req.Valid() {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserDataInvalid,
		})
		return
	}

	user.Email = req.Email
	user.Password = req.Password
	userID := getUserID(user)
	user.UserID = userID

	// User is not registered
	if err = userIsAlreadySignedUP(user); err == nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserIsNotRegistered,
		})
		return
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

	authorizedUsers = append(authorizedUsers, user)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != "POST" {
		return
	}

	// Is user authorized?
	c, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
		return
	}

	// Need to remove user from authorized list
	var id uint64
	for k, v := range activeSessions {
		if v.token == c.Value {
			id = k
		}
	}

	var idx int
	for i, user := range authorizedUsers {
		if user.UserID == id {
			idx = i
		}
	}

	delete(activeSessions, id);
	authorizedUsers = append(authorizedUsers[:idx], authorizedUsers[idx+1:]...)

	// Set cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	sendLogOutResponse(w, response.LogOutResponse{
		Message: "Logout successfull",
	})
}

func sendLogOutResponse(w http.ResponseWriter, lr response.LogOutResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(lr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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
