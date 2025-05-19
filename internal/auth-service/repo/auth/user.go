package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	authEntity "retarget/internal/auth-service/entity/auth"
	optiLog "retarget/pkg/utils/optiLog"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AuthRepositoryInterface interface {
	GetUserByID(id int, requestID string) (*authEntity.User, error)
	GetUserByEmail(email string, requestID string) (*authEntity.User, error)
	GetUserByUsername(username string, requestID string) (*authEntity.User, error)
	CheckEmailOrUsernameExists(email, username string, requestID string) (*authEntity.User, error)
	CreateNewUser(user *authEntity.User, requestID string) error
	CloseConnection() error
}

type AuthRepository struct {
	db          *sql.DB
	asyncLogger *optiLog.AsyncLogger
}

func NewAuthRepository(endPoint string, logger *zap.SugaredLogger) *AuthRepository {
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		logger.Fatalf("Failed to connect to DB: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	asyncLogger := optiLog.NewAsyncLogger(logger, 1000, 100_000)

	return &AuthRepository{
		db:          db,
		asyncLogger: asyncLogger,
	}
}

func withRetry(fn func() error) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if !isTransientError(err) {
			return err
		}

		if attempt < 2 {
			delay := time.Millisecond * time.Duration(100*(1<<uint(attempt)))
			time.Sleep(delay)
		}
	}
	return err
}

func isTransientError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "network error") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "deadlock") ||
		strings.Contains(err.Error(), "try again later") ||
		strings.Contains(err.Error(), "EOF") ||
		strings.Contains(err.Error(), "unexpected EOF") ||
		strings.Contains(err.Error(), "server closed the connection unexpectedly") {
		return true
	}

	return false
}

func (r *AuthRepository) getUserByColumn(columnName, value string, requestID string) (*authEntity.User, error) {
	var user *authEntity.User
	var err error

	retryFn := func() error {
		user, err = r.getUserByColumnOnce(columnName, value, requestID)
		return err
	}

	err = withRetry(retryFn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) getUserByColumnOnce(columnName, value string, requestID string) (*authEntity.User, error) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT id, username, email, password, description, balance, role FROM auth_user WHERE %s = $1", columnName)

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Executing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query":  query,
			"column": columnName,
			"value":  value,
		}))

	row := r.db.QueryRowContext(ctx, query, value)
	user := &authEntity.User{}

	scanErr := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Description,
		&user.Balance,
		&user.Role,
	)

	duration := time.Since(startTime).Milliseconds()

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User not found",
				optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
					columnName: value,
				}))
			return nil, scanErr
		}
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "Database query failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"error": scanErr.Error(),
			}))
		return nil, fmt.Errorf("database error: %w", scanErr)
	}

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User retrieved successfully",
		optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
			"userID":   user.ID,
			"username": user.Username,
		}))

	return user, nil
}

func (r *AuthRepository) GetUserByID(id int, requestID string) (*authEntity.User, error) {
	var user *authEntity.User
	var err error

	retryFn := func() error {
		user, err = r.getUserByIDOnce(id, requestID)
		return err
	}

	err = withRetry(retryFn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) getUserByIDOnce(id int, requestID string) (*authEntity.User, error) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, username, email, password, description, balance, role FROM auth_user WHERE id = $1"
	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Executing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query": query,
			"id":    id,
		}))

	row := r.db.QueryRowContext(ctx, query, id)
	user := &authEntity.User{}

	scanErr := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Description,
		&user.Balance,
		&user.Role,
	)

	duration := time.Since(startTime).Milliseconds()

	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User not found",
				optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
					"id": id,
				}))
			return nil, scanErr
		}
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "Database query failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"error": scanErr.Error(),
			}))
		return nil, fmt.Errorf("database error: %w", scanErr)
	}

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User retrieved successfully",
		optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
			"userID":   user.ID,
			"username": user.Username,
		}))

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(email string, requestID string) (*authEntity.User, error) {
	return r.getUserByColumn("email", email, requestID)
}

func (r *AuthRepository) GetUserByUsername(username string, requestID string) (*authEntity.User, error) {
	return r.getUserByColumn("username", username, requestID)
}

func (r *AuthRepository) CreateNewUser(user *authEntity.User, requestID string) error {
	startTime := time.Now()

	if err := authEntity.ValidateUser(user); err != nil {
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "User validation failed",
			optiLog.MakeLogFields(requestID, time.Since(startTime).Milliseconds(), map[string]interface{}{
				"username": user.Username,
				"error":    err.Error(),
			}))
		return fmt.Errorf("validation error: %w", err)
	}

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Starting user creation",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"username": user.Username,
			"role":     user.Role,
		}))

	var err error
	retryFn := func() error {
		err = r.createNewUserOnce(user, requestID)
		return err
	}

	err = withRetry(retryFn)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) createNewUserOnce(user *authEntity.User, requestID string) error {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO auth_user (username, email, password, description, balance, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Preparing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query": query,
		}))

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Description,
		user.Balance,
		user.Role,
	).Scan(&id)

	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "User creation failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"username": user.Username,
				"error":    err.Error(),
			}))
		return fmt.Errorf("database error: %w", err)
	}

	user.ID = int(id)

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User created successfully",
		optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
			"userID":   user.ID,
			"username": user.Username,
		}))

	return nil
}

func (r *AuthRepository) CheckEmailOrUsernameExists(email, username string, requestID string) (*authEntity.User, error) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT id, username, email, password, description, balance, role 
        FROM auth_user 
        WHERE email = $1 OR username = $2`

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Executing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query":    query,
			"email":    email,
			"username": username,
		}))

	row := r.db.QueryRowContext(ctx, query, email, username)
	user := &authEntity.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Description,
		&user.Balance,
		&user.Role,
	)

	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		if err == sql.ErrNoRows {
			r.asyncLogger.Log(zapcore.DebugLevel, requestID, "No user found with given email or username",
				optiLog.MakeLogFields(requestID, duration, nil))
			return nil, nil
		}
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "Database query failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"error": err.Error(),
			}))
		return nil, fmt.Errorf("database error: %w", err)
	}

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Matching user found",
		optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
			"userID":   user.ID,
			"username": user.Username,
			"email":    user.Email,
		}))

	return user, nil
}

func (r *AuthRepository) CloseConnection() error {
	r.asyncLogger.Close()
	return r.db.Close()
}
