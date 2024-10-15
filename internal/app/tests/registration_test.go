package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	internal_errors "pinset/internal/errors"
	"pinset/internal/models"
	"pinset/internal/models/request"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserIsAlreadySignedUP(t *testing.T) {
	testData := map[string]struct {
		input models.User
		err   error
	}{
		"Not registered": {
			input: models.User{
				UserID:       2,
				UserName:     "admin1",
				NickName:     "admin1",
				Email:        "example1@test.com",
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
			err: nil,
		},
		"Already registered": {
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
			err: internal_errors.ErrUserAlreadyRegistered,
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			resultErr := handlers.TestableUserIsAlreadySignedUP(test.input)

			if resultErr != nil {
				if !errors.Is(resultErr, test.err) {
					t.Errorf("expected error %v , got %v", test.err, resultErr)
				}
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	testData := map[string]struct {
		input  models.User
		result uint64
	}{
		"User exists": {
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
			result: 1,
		},
		"User doesnt exist": {
			input:  models.User{},
			result: 0,
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			result := handlers.TestableGetUserID(test.input)

			require.Equal(t, test.result, result, "not equal")
		})
	}
}

func TestSignUp(t *testing.T) {
	type ExpectedReturn struct {
		StatusCode    int
		SessionCookie string
		ErrorMessage  string
	}

	testTable := []struct {
		name           string
		req            request.SignUPRequest
		expectedReturn ExpectedReturn
	}{
		{
			name: "Valid credentials",
			req: request.SignUPRequest{
				UserName: "Leva",
				NickName: "Leva1488",
				Email:    "leva@test.com",
				Password: "12345678Q",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid credentials",
			req: request.SignUPRequest{
				Email:    "invalid@example.com",
				Password: "wrongpassword",
			},
			expectedReturn: ExpectedReturn{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: internal_errors.ErrUserDataInvalid.Error(),
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(testCase.req)
			req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handlers.SignUp(rr, req)

			require.Equal(t, testCase.expectedReturn.StatusCode, rr.Code, "not equal")
		})
	}
}
