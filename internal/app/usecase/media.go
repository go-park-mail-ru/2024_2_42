package usecase

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"

	internal_errors "pinset/internal/errors"
)

func NewMediaUsecase(repo MediaRepository, userRepo UserRepository) delivery.MediaUsecase {
	return &MediaUsecaseController{
		repo:     repo,
		userRepo: userRepo,
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
	pinSet, err := muc.repo.GetAllPins(userID)
	if err != nil {
		return nil, fmt.Errorf("feed usecase: %w", err)
	}

	for _, pin := range pinSet {
		pin.AuthorInfo, err = muc.GetPinAuthorNickNameByUserID(pin.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("feed usecase GetPinAuthorNickNameByUserID: %w", err)
		}

		pin.AuthorInfo.FollowingsCount, err = muc.userRepo.GetFollowingsCount(pin.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("feed usecase GetFollowingsCount: %w", err)
		}

		if userID != 0 {
			var availableBoards []*models.Board
			availableBoards, err = muc.repo.GetAllBoardsByOwnerID(userID)
			if err != nil {
				return nil, fmt.Errorf("feed usecase GetAllBoardsByOwnerID: %w", err)
			}

			pin.Boards = availableBoards

			var isBookmarked uint64
			isBookmarked, err = muc.repo.GetBookmarkOnUserPin(userID, pin.PinID)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return nil, fmt.Errorf("feed usecase GetBookmarkUserPin: %w", err)
				}
			}

			if isBookmarked != 0 {
				pin.IsBookmarked = true
			}
		}
	}
	return pinSet, nil
}

func (muc *MediaUsecaseController) GetPinPreviewInfo(pinID uint64) (*models.Pin, error) {
	return muc.repo.GetPinPreviewInfoByPinID(pinID)
}

func (muc *MediaUsecaseController) GetPinPageInfo(pinID uint64) (*models.Pin, error) {
	return muc.repo.GetPinPageInfoByPinID(pinID)
}

func (mrc *MediaUsecaseController) GetPinAuthorNickNameByUserID(userID uint64) (*models.UserPin, error) {
	return mrc.repo.GetPinAuthorNickNameByUserID(userID)
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
	err := muc.repo.CreatePinBookmark(bookmark)
	if err != nil {
		return fmt.Errorf("createPinBookmark usecase: %w", err)
	}

	err = muc.repo.UpdateBookmarksCountIncrease(bookmark.PinID)
	if err != nil {
		return fmt.Errorf("updateBookmarksCountIncrease usecase: %w", err)
	}

	return nil
}

func (muc *MediaUsecaseController) DeletePinBookmarkByOwnerIDAndPinID(bookmark models.Bookmark) error {
	err := muc.repo.DeletePinBookmarkByOwnerIDAndPinID(bookmark)
	if err != nil {
		return fmt.Errorf("deletePinBookmarkByOwnerIDAndPinID usecase: %w", err)
	}

	err = muc.repo.UpdateBookmarksCountDecrease(bookmark.PinID)
	if err != nil {
		return fmt.Errorf("UpdatePinBookmarksByPinID usecase: %w", err)
	}
	return nil
}

//////////////////////// BOARDS //////////////////////////

func (muc *MediaUsecaseController) GetAllUserBoards(ownerID uint64, currUserID uint64) ([]*models.Board, error) {
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

func (muc *MediaUsecaseController) GetBoardPins(boardID uint64) ([]*models.Pin, error) {
	PinIDs, err := muc.repo.GetBoardPinsByBoardID(boardID)

	var pins []*models.Pin

	if err != nil {
		return nil, err
	}
	for _, pinID := range PinIDs {
		pin, err := muc.GetPinPageInfo(pinID)
		if err != nil {
			return nil, err
		}
		pins = append(pins, pin)
	}
	return pins, nil
}

func (muc *MediaUsecaseController) AddPinToBoard(boardID uint64, pinID uint64) error {
	return muc.repo.AddPinToBoard(boardID, pinID)
}

func (muc *MediaUsecaseController) UpdateBookmarksCountIncrease(pinID uint64) error {
	return muc.repo.UpdateBookmarksCountIncrease(pinID)
}

func (muc *MediaUsecaseController) UpdateBookmarksCountDecrease(pinID uint64) error {
	return muc.repo.UpdateBookmarksCountDecrease(pinID)
}

func (muc *MediaUsecaseController) DeletePinFromBoard(boardID uint64, pinID uint64) error {
	return muc.repo.DeletePinFromBoardByBoardIDAndPinID(boardID, pinID)
}
