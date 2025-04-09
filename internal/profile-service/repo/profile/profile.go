package profile

import (
	"database/sql"
	"log"
	entityProfile "retarget/internal/profile-service/entity/profile"
)

type ProfileRepositoryInterface interface {
	UpdateProfileByID(userID int, username, description string) error
	GetProfileByID(userID int) (*entityProfile.Profile, error)
	CloseConnection() error
}

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(endPoint string) *ProfileRepository {
	profileRepo := &ProfileRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	profileRepo.db = db
	return profileRepo
}

func (r *ProfileRepository) UpdateProfileByID(userID int, username, description string) error {
	_, err := r.db.Exec("UPDATE user SET username = $1, description = $2 WHERE id = $3",
		username, description, userID)
	return err
}

func (r *ProfileRepository) GetProfileByID(userID int) (*entityProfile.Profile, error) {
	var profile entityProfile.Profile
	err := r.db.QueryRow("SELECT id, username, email, description, balance, role FROM user WHERE id = $1",
		userID).Scan(&profile.ID, &profile.Username, &profile.Email, &profile.Description, &profile.Balance, &profile.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entityProfile.ErrProfileNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) CloseConnection() error {
	return r.db.Close()
}
