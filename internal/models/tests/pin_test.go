package tests

import (
	"errors"
	internal_errors "pinset/internal/errors"
	"pinset/internal/models"
	"testing"
	"time"
)

func TestUserValied(t *testing.T) {
	testData := []struct {
		name   string
		input  models.User
		result error
	}{
		{
			name: "valid data",
			input: models.User{
				UserID:       1,
				UserName:     "admin",
				NickName:     "admin",
				Email:        "example@test.com",
				Password:     "12345678Q",
				BirthTime:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
				Gender:       "table",
				AvatarUrl:    "",
				Followers:    []models.User{},
				Following:    []models.User{},
				Boards:       []models.Board{},
				CreationTime: time.Now(),
				UpdateTime:   time.Now(),
			},
			result: nil,
		},
		{
			name: "Short nickname",
			input: models.User{
				UserID:       1,
				UserName:     "adm",
				NickName:     "admin",
				Email:        "example@test.com",
				Password:     "12345678Q",
				BirthTime:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
				Gender:       "table",
				AvatarUrl:    "",
				Followers:    []models.User{},
				Following:    []models.User{},
				Boards:       []models.Board{},
				CreationTime: time.Now(),
				UpdateTime:   time.Now(),
			},
			result: internal_errors.ErrUserDataInvalid,
		},
		{
			name: "Short password",
			input: models.User{
				UserID:       1,
				UserName:     "adm",
				NickName:     "admin",
				Email:        "example@test.com",
				Password:     "123456",
				BirthTime:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().Location()),
				Gender:       "table",
				AvatarUrl:    "",
				Followers:    []models.User{},
				Following:    []models.User{},
				Boards:       []models.Board{},
				CreationTime: time.Now(),
				UpdateTime:   time.Now(),
			},
			result: internal_errors.ErrUserDataInvalid,
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			resultErr := test.input.Valid()

			if resultErr != nil {
				if !errors.Is(resultErr, test.result) {
					t.Errorf("expected error %v , got %v", test.result, resultErr)
				}
			}
		})
	}
}
