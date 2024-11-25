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
		GetUserIDWithEmail(email string) (uint64, error)
		CreateUser(*models.User) (uint64, error)
		CheckUserByEmail(*models.User) (bool, error)
		GetUserAvatar(uint64) (string, error)
		GetUserInfo(*models.User, uint64) (*models.UserProfile, error)
		GetUserInfoPublic(uint64) (*response.UserProfileResponse, error)
		CheckUserCredentials(*models.User) error
		UpdateUserInfo(*models.User) error
		UpdateUserPassword(*models.User) error
		DeleteUserByID(uint64) error

		GetUsersByParams(*models.UserSearchParams) ([]*models.UserInfo, error)

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
		GetPinAuthorNickNameByUserID(userID uint64) (*models.UserPin, error)
		UpdatePinInfoByPinID(pin *models.Pin) error
		UpdatePinViewsByPinID(pinID uint64) error
		UpdatePinUpdateTimeByPinID() error
		DeletePinByPinID(pinID uint64) error
		GetAllCommentariesByPinID(pinID uint64) ([]*models.Comment, error)
		GetPinBookmarksNumberByPinID(pinID uint64) (uint64, error)
		GetBookmarkOnUserPin(ownerID, pinID uint64) (uint64, error)
		CreatePinBookmark(bookmark *models.Bookmark) error
		DeletePinBookmarkByOwnerIDAndPinID(bookmark models.Bookmark) error
		UpdateBookmarksCountIncrease(pinID uint64) error
		UpdateBookmarksCountDecrease(pinID uint64) error

		GetBoardPinsByBoardID(boardID uint64) ([]uint64, error)
		AddPinToBoard(boardID uint64, pinID uint64) error
		DeletePinFromBoardByBoardIDAndPinID(boardID uint64, pinID uint64) error

		GetAllBoardsByOwnerID(ownerID uint64) ([]*models.Board, error)
		GetBoardByBoardID(boardID uint64) (*models.Board, error)
		CreateBoard(board *models.Board) error
		UpdateBoardByBoardID(board *models.Board) error
		DeleteBoardByBoardID(boardID uint64) error

		GetBucketNameForContentType(fileType string) string
		HasCorrectContentType(string) bool
		UploadMedia(string, string, io.Reader, int64) (string, error)

		CreateChat() (*models.ChatCreateInfo, error)
		AddUserToChat(chatID uint64, userID uint64) error
		GetChatUsers(chatID uint64) ([]uint64, error)
		GetUserChats(userID uint64) ([]uint64, error)
		DeleteChat(chatID uint64) error

		CreateMessage(msg *models.Message) (*models.MessageCreateInfo, error)
		DeleteMessage(messageID uint64) error
		UpdateMessage(msg *models.MessageUpdate) error
		GetChatMessages(chatID uint64) ([]*models.MessageInfo, error)
	}

	UserOnlineRepo interface {
		IsOnlineUser(userID uint64) bool
		GetOnlineUser(userID uint64) *models.ChatUser
		AddOnlineUser(user *models.ChatUser)
		DeleteOnlineUser(userID uint64)
		NumUsersOnline() int
	}
)

// Controllers
type (
	UserUsecaseController struct {
		repo           UserRepository
		mediaRepo      MediaRepository
		authParameters configs.AuthParams
	}

	MediaUsecaseController struct {
		repo     MediaRepository
		userRepo UserRepository
	}

	MessageUsecaseController struct {
		mediaRepo      MediaRepository
		userOnlineRepo UserOnlineRepo
		userRepo       UserRepository
	}
)
