package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type AuthenticatorInterface interface {
	Authenticate(cookie string) (int, int, error)
}

type Authenticator struct {
	redisClient *redis.Client
}

func NewAuthenticator(endPoint string, password string, db string) (*Authenticator, error) {
	db_string, err := strconv.Atoi(db)
	if err != nil {
		return nil, fmt.Errorf("error in config: %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     endPoint,
		Password: password,
		DB:       db_string,
	})

	_, err = redisClient.Ping(context.Background()).Result()
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

	var userID, role int
	_, err = fmt.Sscanf(val, "%d:%d", &userID, &role)
	if err != nil {
		return 0, -1, fmt.Errorf("error parse of value: %w", err)
	}

	return userID, role, nil
}
