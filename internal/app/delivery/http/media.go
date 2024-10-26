package delivery

import (
	"mime"
	"net/http"
	"pinset/pkg/utils"
	"strings"

	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"
)

const (
	successfullUploadMessage = "file(s) successfully uploaded"
)

func (mdc *MediaDeliveryController) GetMedia(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", w.FormDataContentType())
	// SendMediaResponse()
}

func (mdc *MediaDeliveryController) UploadMedia(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(contentType, "multipart/") {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	err = r.ParseMultipartForm(32 * utils.MiB)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	// fileHeaders are accessible only after ParseMultipartForm is called
	fileChunks := r.MultipartForm.File["file"]
	err = mdc.Usecase.UploadMedia(fileChunks)
	if err != nil {
		internal_errors.SendErrorResponse(w, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	SendMediaUploadResponse(w, response.MediaUploadResponse{
		Message: successfullUploadMessage,
	})
}
