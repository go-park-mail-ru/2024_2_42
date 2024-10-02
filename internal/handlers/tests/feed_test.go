package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"youpin/internal/handlers" // Импортируем пакет, где определена функция Feed
	"youpin/internal/models"   // Импортируем пакет, где определены модели данных

	"github.com/stretchr/testify/assert"
)

// Тестируем успешный вызов Feed с методом GET
func TestFeedGETSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Feed)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var feed models.Feed
	err = json.NewDecoder(rr.Body).Decode(&feed)
	assert.NoError(t, err)

	// Проверяем, что в ленте есть данные
	assert.Greater(t, len(feed.Pins), 0)
}

// Тестируем вызов Feed с методом OPTIONS
func TestFeedOPTIONS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Feed)

	handler.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул 200 OK для OPTIONS запроса
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Тестируем вызов Feed с не GET методом (например, POST)
func TestFeedNonGETMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Feed)

	handler.ServeHTTP(rr, req)

	// Ожидаем, что сервер не поддержит этот метод и вернет соответствующий текст
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "For now only GET method is allowed")
}
