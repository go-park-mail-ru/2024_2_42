package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"youpin/internal/handlers"
	"youpin/internal/models/request"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	type ExpectedReturn struct {
		StatusCode    int
		SessionCookie string
		ErrorMessage  string
	}

	testData := []struct {
		name           string
		req            request.LoginRequest
		expectedReturn ExpectedReturn
	}{
		{
			name: "Valid credentials",
			req: request.LoginRequest{
				Email:    "example@test.com",
				Password: "12345678Q",
			},
			expectedReturn: ExpectedReturn{
				StatusCode:    http.StatusOK,
				SessionCookie: "session_token",
			},
		},
		{
			name: "Invalid credentials",
			req: request.LoginRequest{
				Email:    "example1@test.com",
				Password: "12345678Q",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusForbidden,
			},
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(testCase.req)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handlers.LogIn(rr, req)

			require.Equal(t, testCase.expectedReturn.StatusCode, rr.Code, "not equal")
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	type ExpectedReturn struct {
		StatusCode   int
		ErrorMessage string
	}

	testData := []struct {
		name           string
		req            request.LoginRequest
		expectedReturn ExpectedReturn
	}{
		{
			name: "Valid credentials",
			req: request.LoginRequest{
				Email:    "example@test.com",
				Password: "12345678Q",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid credentials",
			req: request.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongPassword",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusForbidden,
			},
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			// Отправка запроса на логин
			reqBody, _ := json.Marshal(testCase.req)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handlers.LogIn(rr, req)

			// Проверка куки
			cookie, _ := rr.Cookie("session_cookie")

			reqIsAuth, err := http.NewRequest(http.MethodGet, "/is_authorized", nil)
			if err != nil {
				t.Fatal(err)
			}
			reqIsAuth.Header.Set("Content-Type", "application/json")
			if cookie != nil {
				reqIsAuth.AddCookie(cookie)
			}

			rrIsAuth := httptest.NewRecorder()
			handlers.IsAuthorized(rrIsAuth, reqIsAuth)
			fmt.Println(reqIsAuth.Cookie("session_token"))
			require.Equal(t, testCase.expectedReturn.StatusCode, rrIsAuth.Code, "not equal")
			fmt.Println(rrIsAuth.Result().StatusCode)
		})
	}
}
