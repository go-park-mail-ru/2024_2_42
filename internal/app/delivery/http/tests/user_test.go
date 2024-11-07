package tests

import (
	"fmt"
	"pinset/internal/app/models"
	"pinset/internal/app/models/request"
	"pinset/internal/app/models/response"
	"pinset/internal/app/usecase"
	"pinset/internal/app/usecase/mock_usecase"
	internal_errors "pinset/internal/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogIn(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User, req request.LoginRequest)

	testCases := []struct {
		name          string
		inputUser     *models.User
		req           request.LoginRequest
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "User exists",
			inputUser: &models.User{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			req: request.LoginRequest{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
				Token:    "12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User, req request.LoginRequest) {
				s.CheckUserByEmail(user).Return(true, nil)
				s.CheckUserCredentials(user).Return(nil)
				s.UserHasActiveSession(req.Token).Return(false)
				s.GetLastUserID().Return(uint64(1), fmt.Errorf(""))
			},
			expectedError: fmt.Errorf("logIn getLastUserID: "),
		},
		{
			name: "User does not exist",
			inputUser: &models.User{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			req: request.LoginRequest{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User, req request.LoginRequest) {
				s.CheckUserByEmail(user).Return(false, nil)
			},
			expectedError: internal_errors.ErrUserDoesntExists,
		},
		{
			name: "Error from CheckUserByEmail",
			inputUser: &models.User{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			req: request.LoginRequest{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User, req request.LoginRequest) {
				s.CheckUserByEmail(user).Return(false, fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("signUp after UserAlreadySignedUp: db error"),
		},
		{
			name: "Error from checkCredentials",
			inputUser: &models.User{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			req: request.LoginRequest{
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User, req request.LoginRequest) {
				s.CheckUserByEmail(user).Return(true, nil)
				s.CheckUserCredentials(user).Return(fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("login checkCredentials: db error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			login := mock_usecase.NewMockUserRepository(c)
			testCase.mockBehavior(login.EXPECT(), testCase.inputUser, testCase.req)

			userUsecase := usecase.NewUserUsecase(login)

			req := testCase.req

			_, err := userUsecase.LogIn(req)

			if testCase.expectedError != nil {
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("expected error: %v, got %v.", testCase.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User)

	testCases := []struct {
		name          string
		inputUser     *models.User
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "User not signed up",
			inputUser: &models.User{
				UserName: "dima",
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.CheckUserByEmail(user).Return(false, nil)
				s.CreateUser(user).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User already signed up",
			inputUser: &models.User{
				UserName: "dima",
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.CheckUserByEmail(user).Return(false, internal_errors.ErrUserAlreadyExists)
			},
			expectedError: fmt.Errorf("signUp after UserAlreadySignedUp: %s", internal_errors.ErrUserAlreadyExists),
		},
		{
			name: "Error db error",
			inputUser: &models.User{
				UserName: "dima",
				Email:    "dima@dima.ru",
				Password: "Dima12345",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.CheckUserByEmail(user).Return(false, fmt.Errorf("db error"))
			},
			expectedError: fmt.Errorf("signUp after UserAlreadySignedUp: db error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			login := mock_usecase.NewMockUserRepository(c)
			testCase.mockBehavior(login.EXPECT(), testCase.inputUser)

			userUsecase := usecase.NewUserUsecase(login)

			req := testCase.inputUser

			err := userUsecase.SignUp(req)

			if testCase.expectedError != nil {
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("expected error: %v, got %v.", testCase.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}

func TestGetUserInfo(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User)

	testCases := []struct {
		name          string
		user          *models.User
		userProfile   response.UserProfileResponse
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "Success response",
			user: &models.User{
				UserID:   uint64(1),
				NickName: "Dima",
				UserName: "Dima",
				Email:    "dima@dima.ru",
				Password: "Dima",
			},
			userProfile: response.UserProfileResponse{
				UserName:               "Dima",
				NickName:               "Dima",
				Description:            ptr("I love dogs!"),
				NumOfUserFollowings:    4,
				NumOfUserSubscriptions: 0,
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.GetUserInfo(user).Return(response.UserProfileResponse{
					UserName:    "Dima",
					NickName:    "Dima",
					Description: ptr("I love dogs!"),
				}, nil)
				s.GetFollowingsCount(user.UserID).Return(uint64(4), nil)
				s.GetSubsriptionsCount(user.UserID).Return(uint64(0), nil)
			},
			expectedError: nil,
		},
		{
			name: "Bad GetUserInfo response",
			user: &models.User{
				UserID:   uint64(1),
				NickName: "Dima",
				UserName: "Dima",
				Email:    "dima@dima.ru",
				Password: "Dima",
			},
			userProfile: response.UserProfileResponse{},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.GetUserInfo(user).Return(response.UserProfileResponse{},
					fmt.Errorf(""))
			},
			expectedError: fmt.Errorf("userProfile GetUserInfo usecase: "),
		},
		{
			name: "Bad GetFollowings count",
			user: &models.User{
				UserID:   uint64(1),
				NickName: "Dima",
				UserName: "Dima",
				Email:    "dima@dima.ru",
				Password: "Dima",
			},
			userProfile: response.UserProfileResponse{},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, user *models.User) {
				s.GetUserInfo(user).Return(response.UserProfileResponse{
					UserName:    "Dima",
					NickName:    "Dima",
					Description: ptr("I love dogs!"),
				}, nil)
				s.GetFollowingsCount(user.UserID).Return(uint64(0),
					fmt.Errorf(""))
			},
			expectedError: fmt.Errorf("userProfile GetFollowingsCount usecase: "),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := mock_usecase.NewMockUserRepository(c)
			testCase.mockBehavior(mockRepo.EXPECT(), testCase.user)

			userUsecase := usecase.NewUserUsecase(mockRepo)

			res, err := userUsecase.GetUserInfo(testCase.user)

			if testCase.expectedError != nil {
				if err == nil || err.Error() != testCase.expectedError.Error() {
					t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
				}
			} else {
				assert.Equal(t, res, testCase.userProfile)
			}
		})
	}
}

func TestUpdateUserInfo(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockUserRepositoryMockRecorder, token string, user *models.User)

	testCases := []struct {
		name          string
		token         string
		user          *models.User
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:  "Unsuccessful update.",
			token: "12345",
			user: &models.User{
				UserID:   uint64(1),
				NickName: "Dima",
				UserName: "Dima",
				Email:    "dima@dima.ru",
				Password: "Dima",
			},
			mockBehavior: func(s *mock_usecase.MockUserRepositoryMockRecorder, token string, user *models.User) {
				s.UpdateUserInfo(user).Return(fmt.Errorf(""))
			},
			expectedError: fmt.Errorf("updateUserPasswordByID isAuthorized: %s", "token contains an invalid number of segments"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := mock_usecase.NewMockUserRepository(c)
			testCase.mockBehavior(mockRepo.EXPECT(), testCase.token, testCase.user)

			userUsecase := usecase.NewUserUsecase(mockRepo)

			err := userUsecase.UpdateUserInfo(testCase.token, testCase.user)

			if testCase.expectedError != nil {
				if err == nil || err.Error() != testCase.expectedError.Error() {
					t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
