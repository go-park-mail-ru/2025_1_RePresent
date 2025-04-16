package repo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	authEntity "retarget/internal/auth-service/entity/auth"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AuthRepositoryInterface interface {
	GetUserByID(id int) (*authEntity.User, error)
	GetUserByEmail(email string) (*authEntity.User, error)
	GetUserByUsername(email string) (*authEntity.User, error)
	CreateNewUser(user *authEntity.User) error

	CloseConnection() error
}

type AuthRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
} // TODO: Переделать коннект в эндпойнт

func NewAuthRepository(endPoint string, logger *zap.SugaredLogger) *AuthRepository {
	userRepo := &AuthRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	userRepo.db = db
	userRepo.logger = logger
	return userRepo
}

func (r *AuthRepository) GetUserByID(id int) (*authEntity.User, error) {
	startTime := time.Now()
	query := "SELECT id, username, email, password, description, balance, role FROM auth_user WHERE id = $1"

	r.logger.Debugw("Executing SQL query",
		"query", query,
		"id", id,
	)
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
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugw("User not found",
				"id", id,
				"timeTakenMs", time.Since(startTime).Milliseconds(),
			)
			return nil, fmt.Errorf("user not found")
		}
		r.logger.Debugw("Database query failed",
			"id", id,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return nil, fmt.Errorf("database error: %w", err)
	}
	r.logger.Debugw("User retrieved successfully",
		"userID", user.ID,
		"username", user.Username,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*authEntity.User, error) {
	startTime := time.Now()
	query := "SELECT id, username, email, password, description, balance, role FROM auth_user WHERE email = $1"

	r.logger.Debugw("Executing SQL query",
		"query", query,
		"email", email,
	)

	row := r.db.QueryRow(query, email)
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

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugw("User not found",
				"email", email,
				"timeTakenMs", time.Since(startTime).Milliseconds(),
			)
			return nil, fmt.Errorf("user not found")
		}

		r.logger.Debugw("Database query failed",
			"email", email,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return nil, fmt.Errorf("database error: %w", err)
	}

	r.logger.Debugw("User retrieved successfully",
		"userID", user.ID,
		"username", user.Username,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)

	return user, nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*authEntity.User, error) {
	startTime := time.Now()
	query := "SELECT id, username, email, password, description, balance, role FROM auth_user WHERE username = $1"

	r.logger.Debugw("Executing SQL query",
		"query", query,
		"username", username,
	)

	row := r.db.QueryRow(query, username)
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

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugw("User not found",
				"username", username,
				"timeTakenMs", time.Since(startTime).Milliseconds(),
			)
			return nil, fmt.Errorf("user not found")
		}

		r.logger.Debugw("Database query failed",
			"username", username,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return nil, fmt.Errorf("database error: %w", err)
	}

	r.logger.Debugw("User retrieved successfully",
		"userID", user.ID,
		"username", user.Username,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)

	return user, nil
}

func (r *AuthRepository) CreateNewUser(user *authEntity.User) error {
	startTime := time.Now()

	// Логируем начало операции (без sensitive данных)
	r.logger.Debugw("Starting user creation",
		"username", user.Username,
		"role", user.Role,
	)

	err := authEntity.ValidateUser(user)
	if err != nil {
		r.logger.Debugw("User validation failed",
			"username", user.Username,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("validation error: %w", err)
	}

	query := "INSERT INTO auth_user (username, email, password, description, balance, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	r.logger.Debugw("Preparing SQL query",
		"query", query,
	)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		r.logger.Debugw("Prepare statement failed",
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("database error: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(
		user.Username,
		user.Email,
		user.Password,
		user.Description,
		user.Balance,
		user.Role,
	).Scan(&id)

	if err != nil {
		r.logger.Debugw("User creation failed",
			"username", user.Username,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("database error: %w", err)
	}

	user.ID = int(id)

	r.logger.Debugw("User created successfully",
		"userID", user.ID,
		"username", user.Username,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)

	return nil
}
func (r *AuthRepository) CloseConnection() error {
	return r.db.Close()
}
