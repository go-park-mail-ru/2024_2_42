package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	internal_errors "pinset/internal/errors"
	"pinset/internal/models/response"
)

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
