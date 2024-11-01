package delivery

import (
	"mime"
	"net/http"
	"pinset/pkg/utils"
	"strings"

	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"

	"github.com/sirupsen/logrus"
)

const (
	successfullGetMessage = "file extracted"
	successfullUploadMessage = "file(s) successfully uploaded"
)

func (mdc *MediaDeliveryController) GetMedia(w http.ResponseWriter, r *http.Request) {
	// not implemented
}

func (mdc *MediaDeliveryController) UploadMedia(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(contentType, "multipart/") {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	err = r.ParseMultipartForm(32 * utils.MiB)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	// fileHeaders are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]
	uploadedMediaIds, err := mdc.Usecase.UploadMedia(files)
	if err != nil {
		if internal_errors.IsInternal(err) {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				Internal: err,
			})
		} else {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInternalServerError,
			})
		}
		return
	}

	mdc.Logger.WithFields(logrus.Fields{
		"files": uploadedMediaIds,
	}).Info("Upload successfull")

	SendMediaUploadResponse(w, mdc.Logger, response.MediaUploadResponse{
		Message: successfullUploadMessage,
	})
}
