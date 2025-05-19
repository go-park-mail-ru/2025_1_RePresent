package repo

import (
	"database/sql"
	"fmt"
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

	asyncLogger := optiLog.NewAsyncLogger(logger, 1000)

	return &AuthRepository{
		db:          db,
		asyncLogger: asyncLogger,
	}
}

func (r *AuthRepository) getUserByColumn(columnName, value string, requestID string) (*authEntity.User, error) {
	startTime := time.Now()
	query := fmt.Sprintf("SELECT id, username, email, password, description, balance, role FROM auth_user WHERE %s = $1", columnName)

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Executing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query":  query,
			"column": columnName,
			"value":  value,
		}))

	row := r.db.QueryRow(query, value)
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
			r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User not found",
				optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
					columnName: value,
				}))
			return nil, err
		}
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "Database query failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"error": err.Error(),
			}))
		return nil, fmt.Errorf("database error: %w", err)
	}

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User retrieved successfully",
		optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
			"userID":   user.ID,
			"username": user.Username,
		}))

	return user, nil
}

func (r *AuthRepository) GetUserByID(id int, requestID string) (*authEntity.User, error) {
	startTime := time.Now()

	query := "SELECT id, username, email, password, description, balance, role FROM auth_user WHERE id = $1"
	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Executing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query": query,
			"id":    id,
		}))

	row := r.db.QueryRow(query, id)
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
			r.asyncLogger.Log(zapcore.DebugLevel, requestID, "User not found",
				optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
					"id": id,
				}))
			return nil, err
		}
		r.asyncLogger.Log(zapcore.WarnLevel, requestID, "Database query failed",
			optiLog.MakeLogFields(requestID, duration, map[string]interface{}{
				"error": err.Error(),
			}))
		return nil, fmt.Errorf("database error: %w", err)
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

	query := "INSERT INTO auth_user (username, email, password, description, balance, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	r.asyncLogger.Log(zapcore.DebugLevel, requestID, "Preparing SQL query",
		optiLog.MakeLogFields(requestID, 0, map[string]interface{}{
			"query": query,
		}))

	var id int64
	err := r.db.QueryRow(query,
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

	row := r.db.QueryRow(query, email, username)
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
