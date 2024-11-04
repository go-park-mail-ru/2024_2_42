package delivery

import (
	"encoding/json"
	"net/http"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

func sendLogInResponse(w http.ResponseWriter, logger *logrus.Logger, sr response.LogInResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func sendLogOutResponse(w http.ResponseWriter, logger *logrus.Logger, lr response.LogOutResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(lr); err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func SendSignUpResponse(w http.ResponseWriter, logger *logrus.Logger, sr response.SignUpResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func SendIsAuthResponse(w http.ResponseWriter, logger *logrus.Logger, ar response.IsAuthResponse) {
	respJSON, err := json.Marshal(ar)
	if err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(respJSON)
}

func SendMediaUploadResponse(w http.ResponseWriter, logger *logrus.Logger, mur response.MediaUploadResponse) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(mur)
	if err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}
}

func SendUserProfileResponse(w http.ResponseWriter, logger *logrus.Logger, ar response.UserProfileResponse) {
	respJSON, err := json.Marshal(ar)
	if err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(respJSON)
}

func SendUpdateResponse(w http.ResponseWriter, logger *logrus.Logger, sr response.UpdateUserInfo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}
