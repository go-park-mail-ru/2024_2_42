package tests

import (
	"database/sql"

	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
	"pinset/internal/app/session"

	internal_errors "pinset/internal/errors"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type Test struct {
	userID uint64
	user   *models.User
	err    error
}

var (
	exampleUser = &models.User{
		UserID:      1,
		UserName:    "dima",
		NickName:    "dima",
		Email:       "dima@dima.ru",
		Password:    "Dima12345",
		Description: "I love dogs!",
		Gender:      "Male",
		AvatarUrl:   "12345",
	}
)

func strPtr(s string) *string {
	return &s
}

func (m *MockUserRepository) LogIn(req request.LoginRequest) (string, error) {

	return "", nil
}

func (m *MockUserRepository) GetLastUserID() (uint64, error) {
	return exampleUser.UserID, nil
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	return nil
}

func (m *MockUserRepository) CheckUserByEmail(user *models.User) (bool, error) {
	if user.Email == "dima1@dima.ru" {
		return false, internal_errors.ErrUserDoesntExists
	}
	if user.Email == "dima@dima.ru" {
		return true, nil
	}
	return false, sql.ErrNoRows
}

func (m *MockUserRepository) GetUserInfo(user *models.User) (response.UserProfileResponse, error) {
	if user.UserID == 1 {
		return response.UserProfileResponse{
			UserName:    user.UserName,
			NickName:    user.NickName,
			Description: &user.Description,
			Gender:      strPtr("Male"),
		}, nil
	}
	return response.UserProfileResponse{}, sql.ErrNoRows
}

func (m *MockUserRepository) CheckUserCredentials(user *models.User) error {
	return nil
}

func (m *MockUserRepository) UpdateUserInfo(user *models.User) error {
	return nil
}

func (m *MockUserRepository) UpdateUserPassword(user *models.User) error {
	return nil
}

func (m *MockUserRepository) DeleteUserByID(userID uint64) error {
	return nil
}

func (m *MockUserRepository) FollowUser(userID, followerID uint64) error {
	return nil
}

func (m *MockUserRepository) UnfollowUser(userID, followerID uint64) error {
	return nil
}

func (m *MockUserRepository) GetAllFollowings(userID, limit uint64) ([]uint64, error) {
	return []uint64{0}, nil
}

func (m *MockUserRepository) GetAllSubscriptions(userID, limit uint64) ([]uint64, error) {
	return []uint64{0}, nil
}

func (m *MockUserRepository) GetFollowingsCount(userID uint64) (uint64, error) {
	return 0, nil
}

func (m *MockUserRepository) GetSubsriptionsCount(userID uint64) (uint64, error) {
	return 0, nil
}

func (m *MockUserRepository) UserHasActiveSession(token string) bool {
	return false
}

func (m *MockUserRepository) Session() *session.SessionsManager {
	args := m.Called()
	return args.Get(0).(*session.SessionsManager)
}
