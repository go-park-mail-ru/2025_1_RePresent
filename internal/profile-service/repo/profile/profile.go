package profile

import (
	"database/sql"
	"log"
	entityProfile "retarget/internal/profile-service/entity/profile"

	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type ProfileRepositoryInterface interface {
	UpdateProfileByID(userID int, username, description string) error
	GetProfileByID(userID int) (*entityProfile.Profile, error)
	CloseConnection() error
}

type ProfileRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewProfileRepository(endPoint string, logger *zap.SugaredLogger) *ProfileRepository {
	profileRepo := &ProfileRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	profileRepo.db = db
	profileRepo.logger = logger
	return profileRepo
}

func (r *ProfileRepository) UpdateProfileByID(userID int, username, description string) error {
	query := "UPDATE auth_user SET username = $1, description = $2 WHERE id = $3"
	r.logger.Debugw("Executing SQL query GetProfileByID", "query", query, "userID", userID)
	startTime := time.Now()
	_, err := r.db.Exec(query, username, description, userID)
	duration := time.Since(startTime)
	if err != nil {
		r.logger.Infow("SQL Error", "userID", userID, "duration", duration)
		return err
	}
	r.logger.Infow("SQL query executed successfully", "userID", userID, "duration", duration)
	return nil
}

func (r *ProfileRepository) GetProfileByID(userID int) (*entityProfile.Profile, error) {
	query := "SELECT id, username, email, description, balance, role FROM auth_user WHERE id = $1"
	r.logger.Debugw("Executing SQL query GetProfileByID", "query", query, "userID", userID)
	var profile entityProfile.Profile
	startTime := time.Now()
	err := r.db.QueryRow(query,
		userID).Scan(&profile.ID, &profile.Username, &profile.Email, &profile.Description, &profile.Balance, &profile.Role)
	duration := time.Since(startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Infow("Profile not found", "userID", userID, "duration", duration)
			return nil, entityProfile.ErrProfileNotFound
		}
		r.logger.Infow("SQL Error", "userID", userID, "duration", duration)
		return nil, err
	}
	r.logger.Infow("SQL query executed successfully", "userID", userID, "duration", duration)
	return &profile, nil
}

func (r *ProfileRepository) CloseConnection() error {
	return r.db.Close()
}
