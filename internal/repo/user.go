package repo

import (
	"database/sql"

	"RE/internal/entity"

	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
