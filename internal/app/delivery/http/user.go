package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"pinset/configs"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"
	"strconv"
	"time"

	"pinset/internal/app/session"

	"github.com/gorilla/mux"
)

const (
	respSignUpSuccessMesssage = "You successfully signed up!"
	respLogOutSuccessMessage  = "Logout successfull"
)

func (udc *UserDeliveryController) LogIn(w http.ResponseWriter, r *http.Request) {
	// Is user authorized?
	c, _ := r.Cookie(session.SessionTokenCookieKey)

	var req request.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	if !req.Valid() {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserDataInvalid,
		})
		return
	}

	if c != nil {
		req.Token = c.Value
	}

	signedToken, err := udc.Usecase.LogIn(req)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	cookie := &http.Cookie{
		Name:     session.SessionTokenCookieKey,
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(72 * time.Hour),
	}

	http.SetCookie(w, cookie)

	sendLogInResponse(w, udc.Logger, response.LogInResponse{
		SessionCookie: signedToken,
	})
}

func (udc *UserDeliveryController) LogOut(w http.ResponseWriter, r *http.Request) {
	// Is user authorized?
	c, err := r.Cookie(session.SessionTokenCookieKey)
	if errors.Is(err, http.ErrNoCookie) {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
		return
	}

	err = udc.Usecase.LogOut(c.Value)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     session.SessionTokenCookieKey,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	sendLogOutResponse(w, udc.Logger, response.LogOutResponse{
		Message: respLogOutSuccessMessage,
	})
}

func (udc *UserDeliveryController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user *models.User = &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	user.Sanitize()

	signedToken, err := udc.Usecase.SignUp(user)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	cookie := &http.Cookie{
		Name:     session.SessionTokenCookieKey,
		Value:    signedToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(72 * time.Hour),
	}

	http.SetCookie(w, cookie)

	SendSignUpResponse(w, udc.Logger, response.SignUpResponse{
		UserID:        user.UserID,
		SessionCookie: signedToken,
	})
}

func (udc *UserDeliveryController) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	// Is user authorized?
	cookie, err := r.Cookie(session.SessionTokenCookieKey)
	if errors.Is(err, http.ErrNoCookie) {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrUserIsNotAuthorized,
		})
		return
	}

	uid, err := udc.Usecase.IsAuthorized(cookie.Value)
	if err == nil {
		SendIsAuthResponse(w, udc.Logger, response.IsAuthResponse{
			UserID: uid,
		})
	} else {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}
}

func (udc *UserDeliveryController) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	uidStr := mux.Vars(r)["user_id"]
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	currUserID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		currUserID = 0
	}

	userInfo, err := udc.Usecase.GetUserInfo(&models.User{UserID: uid}, currUserID)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}
	SendUserProfileResponse(w, udc.Logger, userInfo)

}

func (udc *UserDeliveryController) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		userID = 0
	}

	var req request.UpdateUserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	if !req.Valid() {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserDataInvalid,
		})
		return
	}

	err = udc.Usecase.UpdateUserInfo(&models.User{
		UserID:      userID,
		UserName:    req.UserName,
		NickName:    req.NickName,
		Description: req.Description,
		BirthTime:   &req.BirthTime,
		Gender:      req.Gender,
		AvatarUrl:   &req.AvatarUrl,
	})

	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendInfoResponse(w, udc.Logger, response.ResponseInfo{
		Message: successfullUpdateMessage,
	})
}

func (udc *UserDeliveryController) GetAvatar(w http.ResponseWriter, r *http.Request) {
	currUserID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserIsNotRegistered,
		})
		return
	}

	userAvatar, err := udc.Usecase.GetUserAvatar(currUserID)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendUserAvatarResponse(w, udc.Logger, response.UserAvatar{
		Message:   "Succes!",
		UserID:    currUserID,
		AvatarUrl: userAvatar,
	})
}

func (udc *UserDeliveryController) GetUsersByParams(w http.ResponseWriter, r *http.Request) {

	currUserID, ok := r.Context().Value(configs.UserIdKey).(uint64)
	if !ok {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			Internal: internal_errors.ErrUserIsNotRegistered,
		})
		return
	}

	var userParams *models.UserSearchParams
	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	res, err := udc.Usecase.GetCompanionsForUser(currUserID, userParams)
	if err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		internal_errors.SendErrorResponse(w, udc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}
}
