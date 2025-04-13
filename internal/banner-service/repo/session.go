package repo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	authEntity "retarget/internal/banner-service/entity"

	uuid "github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionRepositoryInterface interface {
	GetSession(sessionId string) (*authEntity.Session, error)
	AddSession(userId int, role int) (*authEntity.Session, error)
	DelSession(sessionId string) error
	CloseConnection() error

	generateSessionID() (string, error)
}

type SessionRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionRepository(endpoint, password string, db int, ttl time.Duration) *SessionRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     endpoint,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	return &SessionRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *SessionRepository) CloseConnection() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// AddSession creates a new session for the user with random ID
func (r *SessionRepository) AddSession(userId int, role int) (*authEntity.Session, error) {
	ctx := context.Background()

	sessionId, err := r.generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &authEntity.Session{
		ID:        sessionId,
		UserID:    userId,
		Role:      role,
		Expires:   time.Now().Add(r.ttl),
		CreatedAt: time.Now(),
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, sessionId, sessionData, r.ttl).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) GetSession(sessionId string) (*authEntity.Session, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, sessionId).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, authEntity.ErrSessionNotFound
		}
		return nil, err
	}

	var session authEntity.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	if time.Now().After(session.Expires) {
		_ = r.DelSession(sessionId)
		return nil, authEntity.ErrSessionNotFound
	}

	return &session, nil
}

func (r *SessionRepository) DelSession(sessionId string) error {
	ctx := context.Background()
	return r.client.Del(ctx, sessionId).Err()
}

func (r *SessionRepository) generateSessionID() (string, error) {
	sessionId := uuid.NewString()
	return sessionId, nil
}
