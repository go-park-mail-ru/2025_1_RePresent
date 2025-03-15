package repo

import (
	"database/sql"

	"retarget/internal/entity"

	_ "github.com/lib/pq"
)

type UserRepositoryInterface interface {
	GetUserByID(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(email string) (*entity.User, error)
	CreateNewUser(user *entity.User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(id int) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE id = $1", id)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE email = $1", email)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE username = $1", username)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateNewUser(user *entity.User) error {
	err := entity.ValidateUser(user)
	if err != nil {
		return err
	}

	stmt, err := r.db.Prepare("INSERT INTO auth_user (username, email, password, avatar, balance, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(user.Username, user.Email, user.Password, user.Avatar, user.Balance, user.Role).Scan(&id)
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}
