package repository

import (
	"os"
	"sync"
	"time"
)

const SessionTokenExpirationTime = time.Hour * 72

var SECRET = []byte(os.Getenv("JWT_SECRET"))

type SessionsManager struct {
	mu   *sync.RWMutex
	data map[string]uint64
}

func NewSessionManager() *SessionsManager {
	return &SessionsManager{
		mu:   &sync.RWMutex{},
		data: make(map[string]uint64),
	}
}

func (sm *SessionsManager) Create(token string, id uint64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.data[token] = id
}

func (sm *SessionsManager) Remove(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.data, token)
}

func (sm *SessionsManager) Exists(token string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, ok := sm.data[token]

	return ok
}