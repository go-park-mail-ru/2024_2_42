package usecase

import (
	"bytes"
	"mime/multipart"
	"net/http"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"

	internal_errors "pinset/internal/errors"
)

func NewMediaUsecase(repo MediaRepository) delivery.MediaUsecase {
	return &MediaUsecaseController{
		repo: repo,
	}
}

func (muc *MediaUsecaseController) UploadMedia(files []*multipart.FileHeader) ([]string, error) {
	var uploadedMediaUrls []string

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		mediaBytes := make([]byte, fileHeader.Size)
		_, err = file.Read(mediaBytes)
		if err != nil {
			return []string{}, err
		}
		r := bytes.NewReader(mediaBytes)

		// Checking the content type
		fileType := http.DetectContentType(mediaBytes)
		if !muc.repo.HasCorrectContentType(fileType) {
			return []string{}, internal_errors.ErrWrongMediaContentType
		}

		bucketName := muc.repo.GetBucketNameForContentType(fileType)
		uploadedMediaId, err := muc.repo.UploadMedia(bucketName, fileHeader.Filename, r, fileHeader.Size)
		if err != nil {
			return []string{}, err
		}

		uploadedMediaUrls = append(uploadedMediaUrls, uploadedMediaId)
	}

	return uploadedMediaUrls, nil
}

//////////////////////// PINS ////////////////////////////

func (muc *MediaUsecaseController) Feed(userID uint64) ([]*models.Pin, error) {
	return muc.repo.GetAllPins(userID)
}

func (muc *MediaUsecaseController) GetPinPreviewInfo(pinID uint64) (*models.Pin, error) {
	return muc.repo.GetPinPreviewInfoByPinID(pinID)
}

func (muc *MediaUsecaseController) GetPinPageInfo(pinID uint64) (*models.Pin, error) {
	return muc.repo.GetPinPageInfoByPinID(pinID)
}

func (mrc *MediaUsecaseController) GetPinAuthorNameByUserID(userID uint64) (*models.User, error) {
	return mrc.repo.GetPinAuthorNameByUserID(userID)
}

func (muc *MediaUsecaseController) GetPinBookmarksNumber(pinID uint64) (uint64, error) {
	return muc.repo.GetPinBookmarksNumberByPinID(pinID)
}

func (muc *MediaUsecaseController) GetAllCommentaries(pinID uint64) ([]*models.Comment, error) {
	return muc.repo.GetAllCommentariesByPinID(pinID)
}

func (muc *MediaUsecaseController) CreatePin(pin *models.Pin) error {
	return muc.repo.CreatePin(pin)
}

func (muc *MediaUsecaseController) UpdatePinInfo(pin *models.Pin) error {
	return muc.repo.UpdatePinInfoByPinID(pin)
}

func (muc *MediaUsecaseController) UpdatePinViewsNumber(pinID uint64) error {
	return muc.repo.UpdatePinViewsByPinID(pinID)
}

func (muc *MediaUsecaseController) DeletePinByPinID(pinID uint64) error {
	return muc.repo.DeletePinByPinID(pinID)
}

func (muc *MediaUsecaseController) GetBookmarkOnUserPin(ownerID, pinID uint64) (uint64, error) {
	return muc.repo.GetBookmarkOnUserPin(ownerID, pinID)
}

func (muc *MediaUsecaseController) CreatePinBookmark(bookmark *models.Bookmark) error {
	return muc.repo.CreatePinBookmark(bookmark)
}

func (muc *MediaUsecaseController) DeletePinBookmarkByBookmarkID(bookmarkID uint64) error {
	return muc.repo.DeletePinBookmarkByBookmarkID(bookmarkID)
}

//////////////////////// BOARDS //////////////////////////

func (muc *MediaUsecaseController) GetAllUserBoards(ownerID uint64) ([]*models.Board, error) {
	return muc.repo.GetAllBoardsByOwnerID(ownerID)
}

func (muc *MediaUsecaseController) GetBoard(boardID uint64) (*models.Board, error) {
	return muc.repo.GetBoardByBoardID(boardID)
}

func (muc *MediaUsecaseController) CreateBoard(board *models.Board) error {
	return muc.repo.CreateBoard(board)
}

func (muc *MediaUsecaseController) UpdateBoard(board *models.Board) error {
	return muc.repo.UpdateBoardByBoardID(board)
}

func (muc *MediaUsecaseController) DeleteBoard(boardID uint64) error {
	return muc.repo.DeleteBoardByBoardID(boardID)
}
