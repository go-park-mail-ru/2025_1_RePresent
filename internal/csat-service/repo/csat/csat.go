package csat

import (
	"database/sql"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type CsatRepositoryInterface interface {
	AddReview()     // TODO
	GetAllReviews() // TODO

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
