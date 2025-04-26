package csat

import (
	"database/sql"
	"fmt"
	"log"
	csatEntity "retarget/internal/csat-service/entity/csat"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type CsatRepositoryInterface interface {
	AddReview(review csatEntity.Review) (*csatEntity.Review, error)
	GetAllReviews() ([]csatEntity.Review, error)
	GetReviewsByUser(userID int) ([]csatEntity.Review, error)
	GetQuestionsByPage(page string) ([]string, error)

	scanReview(rows *sql.Rows) (*csatEntity.Review, error)
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
func (r *CsatRepository) scanReview(rows *sql.Rows) (*csatEntity.Review, error) {
	var review csatEntity.Review
	var createdAt time.Time
	var id string

	err := rows.Scan(
		&id,
		&review.UserID,
		&review.Question,
		&review.Page,
		&review.Comment,
		&review.Rating,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan review: %w", err)
	}

	review.ID = id
	review.CreatedAt = createdAt
	return &review, nil
}

func (r *CsatRepository) AddReview(review csatEntity.Review) error {
	const addQuery = `
		INSERT INTO reviews (
			user_id, 
			question, 
			page, 
			comment,
			rating
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(addQuery,
		review.UserID,
		review.Question,
		review.Page,
		review.Comment,
		review.Rating,
	)

	if err != nil {
		return fmt.Errorf("failed to insert review: %w", err)
	}

	return nil
}

func (r *CsatRepository) GetAllReviews() ([]csatEntity.Review, error) {
	const query = `
		SELECT 
			id,
			user_id,
			question,
			page,
			comment,
			rating,
			created_at
		FROM reviews
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []csatEntity.Review
	for rows.Next() {
		review, err := r.scanReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, *review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return reviews, nil
}

func (r *CsatRepository) GetReviewsByUser(userID int) ([]csatEntity.Review, error) {
	const query = `
		SELECT 
			id,
			user_id,
			question,
			page,
			comment,
			rating,
			created_at
		FROM reviews
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user reviews: %w", err)
	}
	defer rows.Close()

	var reviews []csatEntity.Review
	for rows.Next() {
		review, err := r.scanReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, *review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return reviews, nil
}

func (r *CsatRepository) GetQuestionsByPage(page string) ([]string, error) {
	const query = `
		SELECT DISTINCT question
		FROM reviews
		WHERE page = ?
		ORDER BY question
	`

	rows, err := r.db.Query(query, page)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []string
	for rows.Next() {
		var question string
		if err := rows.Scan(&question); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, question)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return questions, nil
}
