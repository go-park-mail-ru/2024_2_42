package delivery

import (
	"encoding/json"
	"mime"
	"net/http"
	"pinset/pkg/utils"
	"strconv"
	"strings"

	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	successfullGetMessage         = "file extracted"
	successfullUploadMessage      = "file(s) successfully uploaded"
	successfullPinCreationMessage = "pin successfully created"

	successfullUpdateMessage = "update successfull"

	successfullPinDeletionMessage = "pin deletion successfull"

	successfullBookmarkCreateMessage = "bookmark successfully created"
	successfullBookmarkDeletionMessage = "bookmark successfully deleted"

	successfullBoardDeletion = "board deletion successfull"
)

var lastUploadedMediaUrl string

func (mdc *MediaDeliveryController) UploadMedia(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 * utils.MiB)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	// fileHeaders are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]
	mediaUrls, err := mdc.Usecase.UploadMedia(files)
	lastUploadedMediaUrl = mediaUrls[len(mediaUrls)-1]
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
		"files": mediaUrls,
	}).Info("Upload successfull")

	SendMediaUploadResponse(w, mdc.Logger, response.MediaUploadResponse{
		Message: successfullUploadMessage,
	})
}

////////////////////// PINS //////////////////////

func (mdc *MediaDeliveryController) Feed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Write([]byte("For now only GET method is allowed"))
		return
	}

	feed, err := mdc.Usecase.Feed()
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(feed); err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrFeedNotAccessible,
		})
		return
	}
}

func (mdc *MediaDeliveryController) CreatePin(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	// If media-type entity given then upload it to minio
	if strings.HasPrefix(contentType, "multipart/") {
		mdc.UploadMedia(w, r)
	} else {
		var pin models.Pin
		err := json.NewDecoder(r.Body).Decode(&pin)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
			})
			return
		}

		pin.Sanitize()
		pin.MediaUrl = lastUploadedMediaUrl

		mdc.Logger.WithFields(logrus.Fields{
			"media_url": lastUploadedMediaUrl,
		}).Info("Got media url for new pin")

		err = mdc.Usecase.CreatePin(&pin)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInternalServerError,
			})
			return
		}

		mdc.Logger.WithFields(logrus.Fields{
			"pin_id": pin.PinID,
		}).Info("Pin created successfully")

		SendPinCreatedResponse(w, mdc.Logger, response.PinCreatedResponse{
			PinID:   pin.PinID,
			Message: successfullPinCreationMessage,
		})
	}
}

func (mdc *MediaDeliveryController) GetPinPreview(w http.ResponseWriter, r *http.Request) {
	pinIDStr := mux.Vars(r)["pin_id"]
	userIDStr := mux.Vars(r)["user_id"]
	pinID, err := strconv.ParseUint(pinIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	pin, err := mdc.Usecase.GetPinPreviewInfo(pinID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}
	bookmarksNumber, err := mdc.Usecase.GetPinBookmarksNumber(pinID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}
	author, err := mdc.Usecase.GetPinAuthorNameByUserID(userID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendPinPreviewResponse(w, mdc.Logger, response.PinPreviewResponse{
		AuthorName:            author.UserName,
		AuthorAvatarUrl:       author.AvatarUrl,
		AuthorFollowersNumber: 0,
		MediaUrl:              pin.MediaUrl,
		ViewsNumber:           pin.Views,
		BookmarksNumber:       bookmarksNumber,
	})
}

func (mdc *MediaDeliveryController) GetPinPage(w http.ResponseWriter, r *http.Request) {
	pinIDStr := mux.Vars(r)["pin_id"]
	userIDStr := mux.Vars(r)["user_id"]
	pinID, err := strconv.ParseUint(pinIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	pin, err := mdc.Usecase.GetPinPageInfo(pinID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}
	author, err := mdc.Usecase.GetPinAuthorNameByUserID(userID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendPinPageResponse(w, mdc.Logger, response.PinPageResponse{
		AuthorName:            author.UserName,
		AuthorAvatarUrl:       author.AvatarUrl,
		AuthorFollowersNumber: 0,
		MediaUrl:              pin.MediaUrl,
		Title:                 pin.Title,
		Description:           pin.Description,
		RelatedLink:           pin.RelatedLink,
		Geolocation:           pin.Geolocation,
		CreationTime:          pin.CreationTime,
	})
}

func (mdc *MediaDeliveryController) UpdatePin(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	// New media for pin received
	if strings.HasPrefix(contentType, "multipart/") {
		mdc.UploadMedia(w, r)
	} else {
		var req request.UpdatePinRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
			})
			return
		}

		if !req.Valid() {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				Internal: internal_errors.ErrPinDataInvalid,
			})
			return
		}

		err = mdc.Usecase.UpdatePinInfo(&models.Pin{
			PinID:       req.PinID,
			Title:       req.Title,
			Description: req.Description,
			RelatedLink: req.RelatedLink,
			BoardID:     req.BoardID,
			Geolocation: req.Geolocation,
			MediaUrl:    lastUploadedMediaUrl,
		})

		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInternalServerError,
			})
			return
		}

		SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
			Message: successfullUpdateMessage,
		})
	}
}

func (mdc *MediaDeliveryController) DeletePin(w http.ResponseWriter, r *http.Request) {
	pinIDStr := mux.Vars(r)["pin_id"]
	pinID, err := strconv.ParseUint(pinIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	err = mdc.Usecase.DeletePinByPinID(pinID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
		Message: successfullPinDeletionMessage,
	})
}

func (mdc *MediaDeliveryController) GetBookmark(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["owner_id"]
	ownerID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	pinIDStr := mux.Vars(r)["pin_id"]
	pinID, err := strconv.ParseUint(pinIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	bookmarkID, err := mdc.Usecase.GetBookmarkOnUserPin(ownerID, pinID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendResponseBookmarkExists(w, mdc.Logger, response.ResponseBookmarkExists{
		BookmarkID: bookmarkID,
	})
}

func (mdc *MediaDeliveryController) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var bookmark models.Bookmark
	err := json.NewDecoder(r.Body).Decode(&bookmark)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
		})
		return
	}

	err = mdc.Usecase.CreatePinBookmark(&bookmark)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
		Message: successfullBookmarkCreateMessage,
	})
}

func (mdc *MediaDeliveryController) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkIDStr := mux.Vars(r)["bookmark_id"]
	bookmarkID, err := strconv.ParseUint(bookmarkIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	err = mdc.Usecase.DeletePinBookmarkByBookmarkID(bookmarkID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
		Message: successfullBookmarkDeletionMessage,
	})
}

/////////////////// BOARDS ///////////////////

func (mdc *MediaDeliveryController) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	boards, err := mdc.Usecase.GetAllUserBoards(userID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(boards); err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrInternalServerError,
		})
		return
	}
}

func (mdc *MediaDeliveryController) GetBoard(w http.ResponseWriter, r *http.Request) {
	boardIDStr := mux.Vars(r)["board_id"]
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	board, err := mdc.Usecase.GetBoard(boardID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendBoardResponse(w, mdc.Logger, response.BoardResponse{
		BoardID: board.BoardID,
		OwnerID: board.OwnerID,
		Cover: board.Cover,
		Title: board.Name,
		Description: board.Description,
		Public: board.Public,
		CreationTime: board.CreationTime,
		UpdateTime: board.UpdateTime,
	})
}

func (mdc *MediaDeliveryController) CreateBoard(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	// If media-type entity given then upload it to minio
	if strings.HasPrefix(contentType, "multipart/") {
		mdc.UploadMedia(w, r)
	} else {
		var board models.Board
		err := json.NewDecoder(r.Body).Decode(&board)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
			})
			return
		}

		board.Sanitize()
		board.Cover = lastUploadedMediaUrl

		mdc.Logger.WithFields(logrus.Fields{
			"media_url": lastUploadedMediaUrl,
		}).Info("Got media url for new pin")

		err = mdc.Usecase.CreateBoard(&board)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInternalServerError,
			})
			return
		}

		mdc.Logger.WithFields(logrus.Fields{
			"board_id": board.BoardID,
		}).Info("Board created successfully")

		SendBoardCreatedResponse(w, mdc.Logger, response.BoardCreatedResponse{
			BoardID: board.BoardID,
			Message: successfullPinCreationMessage,
		})
	}
}

func (mdc *MediaDeliveryController) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			General: err, Internal: internal_errors.ErrExpectedMultipartContentType,
		})
		return
	}

	// New media for pin received
	if strings.HasPrefix(contentType, "multipart/") {
		mdc.UploadMedia(w, r)
	} else {
		var req request.UpdateBoardRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInvalidOrMissingRequestBody,
			})
			return
		}

		if !req.Valid() {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				Internal: internal_errors.ErrPinDataInvalid,
			})
			return
		}

		err = mdc.Usecase.UpdateBoard(&models.Board{
			BoardID: req.BoardID,
			Cover: req.Cover,
			Name: req.Title,
			Description: req.Description,
			Public: req.Public,
		})

		if err != nil {
			internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
				General: err, Internal: internal_errors.ErrInternalServerError,
			})
			return
		}

		SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
			Message: successfullUpdateMessage,
		})
	}
}

func (mdc *MediaDeliveryController) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	boardIDStr := mux.Vars(r)["board_id"]
	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	err = mdc.Usecase.DeleteBoard(boardID)
	if err != nil {
		internal_errors.SendErrorResponse(w, mdc.Logger, internal_errors.ErrorInfo{
			Internal: err,
		})
		return
	}

	SendInfoResponse(w, mdc.Logger, response.ResponseInfo{
		Message: successfullBoardDeletion,
	})
}
