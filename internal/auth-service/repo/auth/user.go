package repo

import (
	"database/sql"
	"log"

	authEntity "retarget/internal/auth-service/entity/auth"

	_ "github.com/lib/pq"
)

type AuthRepositoryInterface interface {
	GetUserByID(id int) (*authEntity.User, error)
	GetUserByEmail(email string) (*authEntity.User, error)
	GetUserByUsername(email string) (*authEntity.User, error)
	CreateNewUser(user *authEntity.User) error

	CloseConnection() error
}

type AuthRepository struct {
	db *sql.DB
} // TODO: Переделать коннект в эндпойнт

func NewAuthRepository(endPoint string) *AuthRepository {
	userRepo := &AuthRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	userRepo.db = db
	return userRepo
}

func (r *AuthRepository) GetUserByID(id int) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, description, balance, role FROM user WHERE id = $1", id)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Description, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, description, balance, role FROM user WHERE email = $1", email)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Description, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*authEntity.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password, description, balance, role FROM user WHERE username = $1", username)
	user := &authEntity.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Description, &user.Balance, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) CreateNewUser(user *authEntity.User) error {
	err := authEntity.ValidateUser(user)
	if err != nil {
		return err
	}

	stmt, err := r.db.Prepare("INSERT INTO user (username, email, password, description, balance, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(&user.Username, &user.Email, &user.Password, &user.Description, &user.Balance, &user.Role).Scan(&id)
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

func (r *AuthRepository) CloseConnection() error {
	return r.db.Close()
}
