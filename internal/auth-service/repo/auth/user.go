package repo

import (
	"database/sql"
	"fmt"
	"log"

	authEntity "retarget/internal/auth-service/entity/auth"

	_ "github.com/lib/pq"
)

type UserRepositoryInterface interface {
	GetUserByID(id int) (*authEntity.User, error)
	GetUserByEmail(email string) (*authEntity.User, error)
	GetUserByUsername(email string) (*authEntity.User, error)
	CreateNewUser(user *authEntity.User) error

	// OpenConnection(username, password, dbname, host, port, sslmode string) (*sql.DB, error) В данном случае не требуется, но когда-нибудь в целом можно реализовать
	CloseConnection() error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(username, password, dbname, host string, port int, sslmode string) *UserRepository {
	userRepo := &UserRepository{}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		username, password, dbname, host, port, sslmode))
	if err != nil {
		log.Fatal(err)
	}
	userRepo.db = db
	return userRepo
}

func (r *UserRepository) GetUserByID(id int) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE id = $1", id)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE email = $1", email)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, avatar, balance, role FROM auth_user WHERE username = $1", username)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateNewUser(user *authEntity.User) error {
	err := authEntity.ValidateUser(user)
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

func (r *UserRepository) CloseConnection() error {
	return r.db.Close()
}
