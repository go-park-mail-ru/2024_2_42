package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"
	"time"
)

const (
	respSignUpSuccessMesssage = "You successfully signed up!"

	SessionTokenCookieKey = "session_token"
)

func (udc *UserDeliveryController) LogIn(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
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

	signedToken, err := udc.Usecase.LogIn(req)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	cookie := &http.Cookie{
		Name:     SessionTokenCookieKey,
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

func (udc *UserDeliveryController) LogOut(w http.ResponseWriter, r *http.Request) {
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

	err = udc.Usecase.LogOut(c.Value)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

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

func (udc *UserDeliveryController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	user.Sanitize()

	err = udc.Usecase.SignUp(&user)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendSignUpResponse(w, response.SignUpResponse{
		UserId: user.UserID, Message: respSignUpSuccessMesssage,
	})
}

func (udc *UserDeliveryController) IsAuthorized(w http.ResponseWriter, r *http.Request) {
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

	uid, err := udc.Usecase.IsAuthorized(cookie.Value)
	if err == nil {
		SendIsAuthResponse(w, response.IsAuthResponse{
			UserID: uid,
		})
	} else {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			Internal: err,
		})
	}
}
