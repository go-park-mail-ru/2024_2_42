package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"
)

func sendLogInResponse(w http.ResponseWriter, sr response.LogInResponse) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func sendLogOutResponse(w http.ResponseWriter, lr response.LogOutResponse) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(lr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
	}
}

func SendSignUpResponse(w http.ResponseWriter, sr response.SignUpResponse) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrCantProcessFormData,
		})
		return
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

func SendMediaUploadResponse(w http.ResponseWriter, mur response.MediaUploadResponse) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(mur)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}
}

func SendMediaResponse(w http.ResponseWriter, mr response.MediaResponse) {
	w.WriteHeader(http.StatusOK)
}
