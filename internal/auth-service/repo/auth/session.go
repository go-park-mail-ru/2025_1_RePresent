package repo

import (
	"database/sql"
	"log"

	authEntity "retarget/internal/auth-service/entity/auth"

	_ "github.com/lib/pq"
)

type SessionRepositoryInterface interface {
	GetSessionById(id int) (*authEntity.Session, error)
	AddSession(id int, role int) (*authEntity.Session, error)
	DelSession(id int) (*authEntity.User, error)
	CreateNewUser(user *authEntity.User) error

	CloseConnection() error
}

type SessionRepository struct {
	db *sql.DB
} // TODO: Переделать коннект в эндпойнт

func NewSessionRepository(endPoint string) *SessionRepository {
	userRepo := &SessionRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	userRepo.db = db
	return userRepo
}
