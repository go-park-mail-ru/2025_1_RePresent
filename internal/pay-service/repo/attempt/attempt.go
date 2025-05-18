package attempt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type AttemptRepositoryInterface interface {
	ResetAttemptsByUserID(user_id int) error
	IncrementAttemptsByUserID(user_id int) error
	DecrementAttemptsByUserID(userID int) error
	GetAttemptsByUserID(user_id int) (int, error)
	CloseConnection() error
}

type AttemptRepository struct {
	client      *redis.Client
	ttl         time.Duration
	MaxAttempts int
}

func NewAttemptRepository(endpoint, password string, db int, ttl time.Duration, attempts int) *AttemptRepository {
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

	return &AttemptRepository{
		client:      client,
		ttl:         ttl,
		MaxAttempts: attempts,
	}
}

func (r *AttemptRepository) getKey(userID int) string {
	return fmt.Sprintf("attempts:%d", userID)
}

func (r *AttemptRepository) ResetAttemptsByUserID(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Del(ctx, r.getKey(userID)).Err(); err != nil {
		return fmt.Errorf("failed to reset attempts: %w", err)
	}
	return nil
}

func (r *AttemptRepository) IncrementAttemptsByUserID(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := r.getKey(userID)

	result := r.client.Incr(ctx, key)
	if err := result.Err(); err != nil {
		return fmt.Errorf("failed to increment attempts: %w", err)
	}

	if result.Val() == 1 {
		if err := r.client.Expire(ctx, key, r.ttl).Err(); err != nil {
			return fmt.Errorf("failed to set TTL: %w", err)
		}
	}

	return nil
}

func (r *AttemptRepository) DecrementAttemptsByUserID(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := r.getKey(userID)

	// Атомарный Декремент + Удаление на Lua
	script := redis.NewScript(`
        local current = redis.call("DECR", KEYS[1])
        
        -- delete if value <= 0
        if current <= 0 then
            redis.call("DEL", KEYS[1])
            return 0
        end
        
        return current
    `)

	_, err := script.Run(ctx, r.client, []string{key}).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to decrement attempts: %w", err)
	}

	return nil
}

func (r *AttemptRepository) GetAttemptsByUserID(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := r.getKey(userID)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get attempts: %w", err)
	}

	attempts, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid attempts value in Redis: %w", err)
	}

	return attempts, nil
}

func (r *AttemptRepository) CloseConnection() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
