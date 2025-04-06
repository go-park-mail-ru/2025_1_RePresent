package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
}

type AuthenticatorInterface interface {
	Authenticate(cookie string) (int, int, error)
}

type Authenticator struct {
	redisClient *redis.Client
}

func NewAuthenticator(endPoint string, password string, db int) (*Authenticator, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     endPoint,
		Password: password,
		DB:       db,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connect to Redis: %w", err)
	}

	return &Authenticator{redisClient: redisClient}, nil
}

func (a *Authenticator) Authenticate(cookie string) (int, int, error) {
	val, err := a.redisClient.Get(context.Background(), cookie).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, -1, errors.New("session not found")
		}
		return 0, -1, fmt.Errorf("error read from Redis: %w", err)
	}

	var session SessionData
	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return 0, -1, fmt.Errorf("error decoding JSON: %w", err)
	}

	return session.UserID, session.Role, nil
}
