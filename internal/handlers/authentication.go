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

	"github.com/golang-jwt/jwt"
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

var SECRET = []byte(os.Getenv("JWT_SECRET"))

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
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
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrorNotAllowedMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, errors.ErrorWithInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	user, found := findUser(loginRequest.UserName, loginRequest.Password)
	if !found {
		http.Error(w, errors.ErrorWithInvalidUsernameOrPassword.Error(), http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // токен действителен 72 часа
	})

	cryptedToken, err := token.SignedString(SECRET)
	if err != nil {
		http.Error(w, "Error signing the token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    cryptedToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(72 * time.Hour),
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("Login successful"))
}

func findUser(username, password string) (*models.User, bool) {
	regUsrMutex.Lock()
	defer regUsrMutex.Unlock()

	for _, user := range registeredUsers {
		if user.UserName == username && user.Password == password {
			return &user, true
		}
	}
	return nil, false
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
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	fmt.Fprint(w, "Successful logout")
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
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No session token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("no session token")
		}

		return SECRET, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		w.Write([]byte("Successfully authorized as " + username))
		return
	}

	http.Error(w, "Bad request", http.StatusBadRequest)
}
