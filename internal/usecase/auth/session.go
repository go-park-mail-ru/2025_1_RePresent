package auth

import (
	"retarget/internal/entity"
	"sync"
)

var Sessions = make([]entity.Session, 0)
var mu sync.RWMutex

func AddSession(session entity.Session) error {
	mu.Lock()
	defer mu.Unlock()
	Sessions = append(Sessions, session)
	return nil
}

func DelSession(sessionID string) error {
	mu.Lock()
	defer mu.Unlock()
	for i, session := range Sessions {
		if session.ID == sessionID {
			Sessions = append(Sessions[:i], Sessions[i+1:]...)
			return nil
		}
	}
	return entity.ErrSessionNotFound
}

func GetSession(sessionID string) (*entity.Session, error) {
	mu.RLock()
	defer mu.RUnlock()
	for _, session := range Sessions {
		if session.ID == sessionID {
			return &session, nil
		}
	}
	return nil, entity.ErrSessionNotFound
}
