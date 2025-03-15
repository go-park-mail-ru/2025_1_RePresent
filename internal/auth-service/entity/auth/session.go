package entity

import (
	"sync"
	"time"
)

var sessions = make([]Session, 0)
var mu sync.RWMutex

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Expires   time.Time `json:"expires"`
	CreatedAt time.Time `json:"created_at"`
}

var ErrSessionNotFound = &Error{"Session not found"}

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}

func AddSession(session Session) error {
	mu.Lock()
	defer mu.Unlock()
	sessions = append(sessions, session)
	return nil
}

func DelSession(sessionID string) error {
	mu.Lock()
	defer mu.Unlock()
	for i, session := range sessions {
		if session.ID == sessionID {
			sessions = append(sessions[:i], sessions[i+1:]...)
			return nil
		}
	}
	return ErrSessionNotFound
}

func GetSession(sessionID string) (*Session, error) {
	mu.RLock()
	defer mu.RUnlock()
	for _, session := range sessions {
		if session.ID == sessionID {
			return &session, nil
		}
	}
	return nil, ErrSessionNotFound
}
