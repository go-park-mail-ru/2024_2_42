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
		cookie         *http.Cookie
		expectedReturn ExpectedReturn
	}{
		{
			name: "Valid credentials",
			cookie: &http.Cookie{
				Name:  "session_token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgxNDQ0NTYsImxvZ2luIjoiZXhhbXBsZUB0ZXN0LmNvbSIsInVzZXJfaWQiOjF9.83PGao9P9HNzO10f_J-1CMi_7IzWv-iJBHf8JWpk_Oc",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid credentials",
			cookie: &http.Cookie{
				Name:  "session_token",
				Value: "yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgxNDQ0NTYsImxvZ2luIjoiZXhhbXBsZUB0ZXN0LmNvbSIsInVzZXJfaWQiOjF9.83PGao9P9HNzO10f_J-1CMi_7IzWv-iJBHf8JWpk_Oc",
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusForbidden,
			},
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {

			reqIsAuth, err := http.NewRequest(http.MethodGet, "/is_authorized", nil)
			if err != nil {
				t.Fatal(err)
			}
			reqIsAuth.Header.Set("Content-Type", "application/json")
			if testCase.cookie != nil {
				reqIsAuth.AddCookie(testCase.cookie)
			}

			rrIsAuth := httptest.NewRecorder()
			handlers.IsAuthorized(rrIsAuth, reqIsAuth)
			require.Equal(t, testCase.expectedReturn.StatusCode, rrIsAuth.Code, "not equal")
			fmt.Println(rrIsAuth.Result().StatusCode)
		})
	}
}

func TestLogout(t *testing.T) {
	type ExpectedReturn struct {
		StatusCode   int
		ErrorMessage string
	}

	testData := []struct {
		name           string
		cookie         *http.Cookie
		expectedReturn ExpectedReturn
	}{
		{
			name: "Valid credentials",
			cookie: &http.Cookie{
				Name:  "session_token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgxNDQ0NTYsImxvZ2luIjoiZXhhbXBsZUB0ZXN0LmNvbSIsInVzZXJfaWQiOjF9.83PGao9P9HNzO10f_J-1CMi_7IzWv-iJBHf8JWpk_Oc", // Пример корректного токена
			},
			expectedReturn: ExpectedReturn{
				StatusCode: http.StatusOK,
			},
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			reqLogout, err := http.NewRequest(http.MethodPost, "/logout", nil)
			if err != nil {
				t.Fatal(err)
			}
			reqLogout.Header.Set("Content-Type", "application/json")

			if testCase.cookie != nil {
				reqLogout.AddCookie(testCase.cookie)
			}

			rrLogout := httptest.NewRecorder()

			handlers.LogOut(rrLogout, reqLogout)

			require.Equal(t, testCase.expectedReturn.StatusCode, rrLogout.Code, "not equal")
		})
	}
}
