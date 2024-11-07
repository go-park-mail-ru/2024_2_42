package userRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/internal/app/models"
	"pinset/internal/app/models/response"
	internal_errors "pinset/internal/errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	userID        uint64
	expectedError error
	dbError       error
}

var errMockDB = errors.New("mock db error")

func TestGetLastUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var uid uint64 = 1

	logger := &logrus.Logger{}

	rows := sqlmock.NewRows([]string{"userId"}).AddRow(uid)

	tests := []Test{
		{
			userID:        uid,
			expectedError: nil,
			dbError:       nil,
		},
		{
			userID:        0,
			expectedError: errMockDB,
			dbError:       errMockDB,
		},
		{
			userID:        2,
			expectedError: nil,
			dbError:       nil,
		},
	}

	ProfileManager := NewUserRepository(db, logger)
	for casenum, test := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(GetLastUserID))
		if test.dbError != nil {
			query.WillReturnError(test.dbError)
		} else if casenum == 1 {
			query.WillReturnRows(rows)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"userID"}).AddRow(test.userID))
		}

		userID, err := ProfileManager.GetLastUserID()
		assert.Equalf(t, test.userID, userID, "case [%d]: results must match, want %v, have %v", casenum, test.userID, userID)

		if !errors.Is(err, test.expectedError) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{
		UserName: "test_user",
		NickName: "test",
		Email:    "test@example.com",
		Password: "password",
	}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
		userID        uint64
	}{
		{
			name:          "successful creation",
			expectedError: nil,
			userID:        1,
		},
		{
			name:          "database error",
			expectedError: errMockDB,
			dbError:       errMockDB,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(CreateUser)).
			WithArgs(user.UserName, user.NickName, user.Email, user.Password)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"userID"}).AddRow(tt.userID))
		}

		err := repo.CreateUser(user)
		if !errors.Is(err, tt.expectedError) {
			t.Errorf("%s: unexpected error, got %v, want %v", tt.name, err, tt.expectedError)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestCheckUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{
		Email: "test@example.com",
	}

	tests := []struct {
		name          string
		expectedValid bool
		expectedError error
		dbError       error
		foundEmail    sql.NullInt64
	}{
		{
			name:          "user exists",
			expectedValid: true,
			expectedError: nil,
			foundEmail:    sql.NullInt64{Valid: true},
		},
		{
			name:          "user does not exist",
			expectedValid: false,
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(CheckUserByEmail)).
			WithArgs(user.Email)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"foundEmail"}).AddRow(tt.foundEmail))
		}

		valid, err := repo.CheckUserByEmail(user)
		assert.Equal(t, tt.expectedValid, valid, "%s: expected %v, got %v", tt.name, tt.expectedValid, valid)

		if err != nil {
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("%s: unexpected error, got %v, want %v", tt.name, err, tt.expectedError)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestCheckUserCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password",
	}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
		dbPassword    string
	}{
		{
			name:          "correct password",
			expectedError: nil,
			dbPassword:    "password",
		},
		{
			name:          "incorrect password",
			expectedError: internal_errors.ErrBadPassword,
			dbPassword:    "wrong_password",
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(CheckUserCredentials)).
			WithArgs(user.Email)

		if tt.dbError != nil && !errors.Is(tt.dbError, internal_errors.ErrBadPassword) {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(tt.dbPassword))
		}

		err := repo.CheckUserCredentials(user)
		if err != nil {
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("%s: unexpected error, got %v, want %v", tt.name, err, tt.expectedError)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func ptr(s string) *string {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestGetUserInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{UserID: 1, UserName: "test_user", NickName: "test", Description: "test_description", BirthTime: time.Now(), Gender: "male"}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
		userInfo      response.UserProfileResponse
	}{
		{
			name:          "successful fetch",
			expectedError: nil,
			userInfo: response.UserProfileResponse{
				UserName:    "test_user",
				NickName:    "test",
				Description: ptr("test description"),
				Gender:      ptr("male"),
				BirthTime:   ptrTime(time.Now()),
			},
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(GetUserInfoByID)).
			WithArgs(user.UserID)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{
				"userName", "nickName", "description", "birthTime", "gender", "avatarUrl",
			}).AddRow(
				tt.userInfo.UserName, tt.userInfo.NickName, tt.userInfo.Description, tt.userInfo.BirthTime, tt.userInfo.Gender, nil,
			))
		}

		userInfo, err := repo.GetUserInfo(user)
		assert.Equal(t, tt.expectedError, err, "%s: expected %v, got %v", tt.name, tt.expectedError, err)
		assert.Equal(t, tt.userInfo, userInfo, "%s: expected %v, got %v", tt.name, tt.userInfo, userInfo)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestUpdateUserInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{
		UserID:      1,
		UserName:    "test_user",
		NickName:    "test_nick",
		Description: "test_description",
		BirthTime:   time.Now(),
		Gender:      "male",
		AvatarUrl:   "avatar.jpg",
	}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
	}{
		{
			name:          "successful update",
			expectedError: nil,
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(UpdateUserInfoByID)).
			WithArgs(user.UserName, user.NickName, user.Description, user.BirthTime, user.Gender, user.AvatarUrl, user.UserID)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"userID"}).AddRow(user.UserID))
		}

		err := repo.UpdateUserInfo(user)
		assert.Equal(t, tt.expectedError, err, "%s: expected %v, got %v", tt.name, tt.expectedError, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := &models.User{
		UserID: 1,
	}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
	}{
		{
			name:          "successful password update",
			expectedError: nil,
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(UpdateUserPasswordByID)).
			WithArgs(user.UserID)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"userID"}).AddRow(user.UserID))
		}

		err := repo.UpdateUserPassword(user)
		if err != nil {
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("%s: unexpected error, got %v, want %v", tt.name, err, tt.expectedError)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	tests := []struct {
		name          string
		expectedError error
		dbError       error
	}{
		{
			name:          "user deleted successfully",
			expectedError: nil,
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		mockExec := mock.ExpectExec(regexp.QuoteMeta(DeleteUserByID)).
			WithArgs(1)

		// Если тест предполагает ошибку, устанавливаем ее в ожидании `ExpectExec`
		if tt.dbError != nil {
			mockExec.WillReturnError(tt.dbError)
		} else {
			// В случае успешного выполнения возвращаем успешный результат
			mockExec.WillReturnResult(sqlmock.NewResult(0, 1))
		}

		err := repo.DeleteUserByID(1)
		if err != nil {
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("%s: unexpected error, got %v, want %v", tt.name, err, tt.expectedError)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestFollowUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	user := []models.User{
		{
			UserID: 1,
		},
		{
			UserID: 2,
		},
	}

	tests := []struct {
		name          string
		expectedError error
		dbError       error
	}{
		{
			name:          "successful follow",
			expectedError: nil,
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(FollowUser)).
			WithArgs(1, 2)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows(
				[]string{"ownerID", "followerID"}).AddRow(user[0].UserID, user[1].UserID))
		}

		err := repo.FollowUser(1, 2)
		assert.Equal(t, tt.expectedError, err, "%s: expected %v, got %v", tt.name, tt.expectedError, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestUnfollowUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	tests := []struct {
		name          string
		expectedError error
		dbError       error
	}{
		{
			name:          "successful unfollow",
			expectedError: nil,
		},
		{
			name:          "user not found",
			expectedError: internal_errors.ErrUserDoesntExists,
			dbError:       sql.ErrNoRows,
		},
	}

	user := []models.User{
		{
			UserID: 1,
		},
		{
			UserID: 2,
		},
	}

	for _, tt := range tests {
		query := mock.ExpectQuery(regexp.QuoteMeta(UnfollowUser)).
			WithArgs(1, 2)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(sqlmock.NewRows(
				[]string{"ownerID", "followerID"}).AddRow(user[0].UserID, user[1].UserID))
		}

		err := repo.UnfollowUser(1, 2)
		assert.Equal(t, tt.expectedError, err, "%s: expected %v, got %v", tt.name, tt.expectedError, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestGetAllFollowings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	tests := []struct {
		name          string
		expectedError error
		dbError       error
		followers     []uint64
	}{
		{
			name:          "retrieve followings successfully",
			expectedError: nil,
			followers:     []uint64{1, 2, 3},
		},
		{
			name:          "database error",
			expectedError: fmt.Errorf("getAllFollowings: %w", sql.ErrNoRows),
			dbError:       sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		// Создаем строки с одной колонкой для follower_id
		rows := sqlmock.NewRows([]string{"owner_id"})
		for _, follower := range tt.followers {
			rows.AddRow(follower)
		}

		query := mock.ExpectQuery(regexp.QuoteMeta(GetAllFollowings)).
			WithArgs(1)

		if tt.dbError != nil {
			query.WillReturnError(tt.dbError)
		} else {
			query.WillReturnRows(rows)
		}

		// Вызов тестируемой функции
		followers, err := repo.GetAllFollowings(1, 2)
		assert.Equal(t, tt.expectedError, err, "%s: expected %v, got %v", tt.name, tt.expectedError, err)
		assert.Equal(t, tt.followers, followers, "%s: expected %v, got %v", tt.name, tt.followers, followers)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: there were unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestGetAllSubscriptions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	mock.ExpectQuery(regexp.QuoteMeta(GetAllSubscriptions)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"follower_id"}).AddRow(2))

	result, err := repo.GetAllSubscriptions(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{2}, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetFollowingsCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	tests := []struct {
		name          string
		expectedCount uint64
		expectedError error
		queryError    error
	}{
		{
			name:          "successful count",
			expectedCount: 5,
			expectedError: nil,
		},
		{
			name:          "no rows",
			expectedCount: 0,
			queryError:    sql.ErrNoRows,
			expectedError: nil,
		},
		{
			name:          "database error",
			expectedCount: 0,
			queryError:    fmt.Errorf("some database error"),
			expectedError: fmt.Errorf("psql GetFollowingsCount: %w", fmt.Errorf("some database error")),
		},
	}

	for _, tt := range tests {
		// Ожидаем выполнение запроса COUNT
		mockQuery := mock.ExpectQuery(regexp.QuoteMeta(GetFollowingsCount)).
			WithArgs(1)

		if tt.queryError != nil {
			mockQuery.WillReturnError(tt.queryError)
		} else {
			// Возвращаем количество подписчиков
			mockQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.expectedCount))
		}

		count, err := repo.GetFollowingsCount(1)
		if tt.expectedError != nil {
			assert.EqualError(t, err, tt.expectedError.Error(), tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
		assert.Equal(t, tt.expectedCount, count, tt.name)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: unfulfilled expectations: %v", tt.name, err)
		}
	}
}

func TestGetSubscriptionsCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	logger := &logrus.Logger{}
	repo := NewUserRepository(db, logger)

	tests := []struct {
		name          string
		expectedCount uint64
		expectedError error
		queryError    error
	}{
		{
			name:          "successful count",
			expectedCount: 3,
			expectedError: nil,
		},
		{
			name:          "no rows",
			expectedCount: 0,
			queryError:    sql.ErrNoRows,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		mockQuery := mock.ExpectQuery(regexp.QuoteMeta(GetSubsriptionsCount)).
			WithArgs(1)

		if tt.queryError != nil {
			mockQuery.WillReturnError(tt.queryError)
		} else {
			mockQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.expectedCount))
		}

		count, err := repo.GetSubsriptionsCount(1)
		assert.Equal(t, tt.expectedError, err)
		assert.Equal(t, tt.expectedCount, count)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("%s: unfulfilled expectations: %v", tt.name, err)
		}
	}
}
