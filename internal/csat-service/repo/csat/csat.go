package csat

import (
	"database/sql"
	"fmt"
	"log"
	csatEntity "retarget/internal/csat-service/entity/csat"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type CsatRepositoryInterface interface {
	AddReview(review csatEntity.Review) error                  // TODO
	GetAllReviews() ([]csatEntity.Review, error)               // TODO
	GetReviewsByUser(user_id int) ([]csatEntity.Review, error) // TODO
	GetQuestionsByPage(page string) ([]string, error)          // TODO

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

func (r *CsatRepository) AddReview(review csatEntity.Review) (*csatEntity.Review, error) {
	query := `
        INSERT INTO reviews (
            question, 
            page, 
            comment, 
            rating,
            user_id,       // если есть в структуре
            created_at     // если нужно
        ) VALUES (?, ?, ?, ?, ?, ?)
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		review.Question,
		review.Page,
		review.Comment,
		review.Rating,
		review.User_id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert review: %w", err)
	}

	return &review, nil
}
