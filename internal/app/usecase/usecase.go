package usecase

import (
	"io"
	"pinset/configs"
	"pinset/internal/app/models"
	"pinset/internal/app/models/response"
	"pinset/internal/app/session"
)

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go

// Repository interfaces
type (
	UserRepository interface {
		GetLastUserID() (uint64, error)
		CreateUser(*models.User) (uint64, error)
		CheckUserByEmail(*models.User) (bool, error)
		GetUserInfo(*models.User) (response.UserProfileResponse, error)
		CheckUserCredentials(*models.User) error
		UpdateUserInfo(*models.User) error
		UpdateUserPassword(*models.User) error
		DeleteUserByID(uint64) error

		FollowUser(uint64, uint64) error
		UnfollowUser(uint64, uint64) error
		GetAllFollowings(uint64, uint64) ([]uint64, error)
		GetAllSubscriptions(uint64, uint64) ([]uint64, error)
		GetFollowingsCount(uint64) (uint64, error)
		GetSubsriptionsCount(uint64) (uint64, error)

		UserHasActiveSession(string) bool
		Session() *session.SessionsManager
	}

	MediaRepository interface {
		CreatePin(pin *models.Pin) error
		GetAllPins(uint64) ([]*models.Pin, error)
		GetPinPreviewInfoByPinID(pinID uint64) (*models.Pin, error)
		GetPinPageInfoByPinID(pinID uint64) (*models.Pin, error)
		GetPinAuthorNameByUserID(userID uint64) (*models.User, error)
		UpdatePinInfoByPinID(pin *models.Pin) error
		UpdatePinViewsByPinID(pinID uint64) error
		UpdatePinUpdateTimeByPinID() error
		DeletePinByPinID(pinID uint64) error
		GetAllCommentariesByPinID(pinID uint64) ([]*models.Comment, error)
		GetPinBookmarksNumberByPinID(pinID uint64) (uint64, error)
		GetBookmarkOnUserPin(ownerID, pinID uint64) (uint64, error)
		CreatePinBookmark(bookmark *models.Bookmark) error
		DeletePinBookmarkByBookmarkID(bookmarkID uint64) error

		GetBoardPinsByBoardID(boardID uint64) ([]uint64, error)
		AddPinToBoard(boardID uint64, pinID uint64) error

		GetAllBoardsByOwnerID(ownerID uint64) ([]*models.Board, error)
		GetBoardByBoardID(boardID uint64) (*models.Board, error)
		CreateBoard(board *models.Board) error
		UpdateBoardByBoardID(board *models.Board) error
		DeleteBoardByBoardID(boardID uint64) error

		GetBucketNameForContentType(fileType string) string
		HasCorrectContentType(string) bool
		UploadMedia(string, string, io.Reader, int64) (string, error)
	}
)

// Controllers
type (
	UserUsecaseController struct {
		repo           UserRepository
		authParameters configs.AuthParams
	}

	MediaUsecaseController struct {
		repo MediaRepository
	}
)
