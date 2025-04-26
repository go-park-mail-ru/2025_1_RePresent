package csat

import (
	"database/sql"
	"log"
	csatEntity "retarget/internal/csat-service/entity/csat"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type CsatRepositoryInterface interface {
	AddReview(csatEntity.Review) error                         // TODO
	GetAllReviews() ([]csatEntity.Review, error)               // TODO
	GetReviewsByUser(user_id int) ([]csatEntity.Review, error) // TODO
	GetReviewsByPage(page string) ([]csatEntity.Review, error) // TODO

	CloseConnection() error
}

type CsatRepository struct {
	db *sql.DB
}

func NewCsatRepository(dsn string) *CsatRepository {
	csatRepo := &CsatRepository{}
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal(err)
	}
	csatRepo.db = db
	return csatRepo
}

func (r *CsatRepository) CloseConnection() error {
	return r.db.Close()
}

func (r *CsatRepository) AddReview(csatEntity.Review) (*authEntity.User, error) {}
