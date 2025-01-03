package delivery

import (
	"mime/multipart"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Usecase interfaces
type (
	UserUsecase interface {
		LogIn(request.LoginRequest) (string, error)
		LogOut(string) error
		SignUp(user *models.User) (string, error)
		IsAuthorized(string) (uint64, error)
		GetUserAvatar(uint64) (string, error)
		GetUserInfo(*models.User, uint64) (*models.UserProfile, error)
		GetUserInfoPublic(uint64) (*response.UserProfileResponse, error)
		UpdateUserInfo(*models.User) error
		GetUsersByParams(*models.UserSearchParams) ([]*models.UserInfo, error)
		GetCompanionsForUser(uint64, *models.UserSearchParams) ([]*models.UserInfo, error)
	}

	MediaUsecase interface {
		UploadMedia(files []*multipart.FileHeader) ([]string, error)

		Feed(uint64) ([]*models.Pin, error)
		GetPinPreviewInfo(pinID uint64) (*models.Pin, error)
		GetPinPageInfo(pinID uint64) (*models.Pin, error)
		GetPinAuthorNickNameByUserID(userID uint64) (*models.UserPin, error)
		GetAllCommentaries(pinID uint64) ([]*models.Comment, error)
		CreatePin(pin *models.Pin) error
		UpdatePinInfo(pin *models.Pin) error
		UpdatePinViewsNumber(pinID uint64) error
		DeletePinByPinID(pinID uint64) error

		GetBoardPins(boardID uint64) ([]*models.Pin, error)
		AddPinToBoard(boardID uint64, pinID uint64) error
		DeletePinFromBoard(boardID uint64, pinID uint64) error

		GetBookmarkOnUserPin(ownerID, pinID uint64) (uint64, error)
		CreatePinBookmark(bookmark *models.Bookmark) error
		GetPinBookmarksNumber(pinID uint64) (uint64, error)
		DeletePinBookmarkByOwnerIDAndPinID(bookmark models.Bookmark) error
		UpdateBookmarksCountIncrease(pinID uint64) error
		UpdateBookmarksCountDecrease(pinID uint64) error

		GetAllUserBoards(ownerID uint64, currUserID uint64) ([]*models.Board, error)
		GetBoard(boardID uint64) (*models.Board, error)
		CreateBoard(board *models.Board) error
		UpdateBoard(board *models.Board) error
		DeleteBoard(boardID uint64) error
	}

	MessageUsecase interface {
		AddOnlineUser(user *models.ChatUser)
		IsOnlineUser(userID uint64) bool
		GetOnlineUser(userID uint64) *models.ChatUser
		DeleteOnlineUser(userID uint64)
		NumUsersOnline() int

		GetChatMessages(chatID uint64) ([]*models.MessageInfo, error)
		AddChatMessage(message *models.Message) (*models.MessageCreateInfo, error)
		GetChatUsers(chatID uint64) ([]uint64, error)
		GetUserChats(userID uint64) ([]*models.ChatInfo, error)

		CreateChat(req *models.ChatCreateRequest) (*models.ChatInfo, error)
	}
)

// Controllers
type (
	UserDeliveryController struct {
		Usecase UserUsecase
		Logger  *logrus.Logger
	}

	MediaDeliveryController struct {
		Usecase MediaUsecase
		Logger  *logrus.Logger
	}

	MessageDelieveryController struct {
		Usecase  MessageUsecase
		Logger   *logrus.Logger
		Upgrader websocket.Upgrader
	}
)
