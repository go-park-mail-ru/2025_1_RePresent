package repo

import (
	"database/sql"
	"errors"
	"log"
	"retarget/internal/banner-service/entity"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type BannerRepositoryInterface interface {
	GetBannerByUserId(id int) ([]*entity.Banner, error)
	CreateNewBanner(banner entity.Banner)
	UpdateBanner(banner entity.Banner)
	GetBannerByID(id int) (*entity.Banner, error)
	DeleteBannerByID(owner, id int) error
}

type BannerRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewBannerRepository(endPoint string, logger *zap.SugaredLogger) *BannerRepository {
	bannerRepo := &BannerRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	bannerRepo.db = db
	bannerRepo.logger = logger
	return bannerRepo
}

func (r *BannerRepository) GetBannersByUserId(id int) ([]entity.Banner, error) {
	query := "SELECT id, owner_id, title, description, content, status, link FROM banner WHERE owner_id = $1 AND deleted = FALSE;"
	r.logger.Debugw("Executing SQL query GetProfileByID", "query", query, "userID", id)
	startTime := time.Now()
	rows, err := r.db.Query(query, id)
	duration := time.Since(startTime)
	if err != nil {
		r.logger.Debugw("SQL Error", "userID", id, "duration", duration, "error", err)
		return nil, err
	}
	defer rows.Close()
	banners := []entity.Banner{}

	for rows.Next() {
		banner := entity.Banner{}
		err := rows.Scan(&banner.ID, &banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Link)
		if err != nil {
			r.logger.Debugw("SQL Error", "userID", id, "duration", duration, "error", err)
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		r.logger.Debugw("SQL Error", "userID", id, "duration", duration, "error", err)
		return nil, err
	}
	r.logger.Debugw("SQL  query executed successfully", "userID", id, "duration", duration, "error", err)
	return banners, nil
}

func (r *BannerRepository) CreateNewBanner(banner entity.Banner) error {
	r.logger.Debugw("Executing SQL query СreateNewBanner", "ownerID", banner.OwnerID, "title", banner.Title, "description", banner.Description, "status", banner.Status, "link", banner.Link)
	stmt, err := r.db.Prepare("INSERT INTO banner (owner_id, title, description, content, status, balance, link) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;")
	startTime := time.Now()

	if err != nil {
		r.logger.Debugw("Error preparing SQL statement", "error", err)
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(banner.OwnerID, banner.Title, banner.Description, banner.Content, banner.Status, 0, banner.Link).Scan(&id)
	if err != nil {
		r.logger.Debugw("Error executing query to create new banner", "error", err)
		return err
	}
	duration := time.Since(startTime)
	r.logger.Debugw("Successfully created new banner", "bannerID", id, "duration", duration)

	return nil

}

func (r *BannerRepository) UpdateBanner(banner entity.Banner) error {
	startTime := time.Now()
	query := "UPDATE banner SET title = $1, description = $2, content = $3, link = $4, status = $5 WHERE id = $6"
	r.logger.Debugw("Starting banner update",
		"bannerID", banner.ID,
		"query", query,
	)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		r.logger.Debugw("Failed to prepare SQL statement for banner update",
			"error", err,
			"bannerID", banner.ID,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(banner.Title, banner.Description, banner.Content, banner.Link, banner.Status, banner.ID)
	if err != nil {
		r.logger.Debugw("Failed to execute banner update",
			"error", err,
			"bannerID", banner.ID,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	r.logger.Debugw("Banner updated successfully",
		"bannerID", banner.ID,
		"timeTaken", time.Since(startTime).String(),
	)
	return nil
}

func (r *BannerRepository) GetBannerByID(id int) (*entity.Banner, error) {
	startTime := time.Now()
	query := `
		SELECT owner_id, title, description, content, balance, link, status
		FROM banner
		WHERE id = $1 AND deleted = FALSE;
		`
	r.logger.Debugw("Starting to fetch banner by ID",
		"bannerID", id, "query", query,
	)

	row := r.db.QueryRow(query, id)

	banner := &entity.Banner{}
	err := row.Scan(
		&banner.OwnerID,
		&banner.Title,
		&banner.Description,
		&banner.Content,
		&banner.Balance,
		&banner.Link,
		&banner.Status,
	)
	if err != nil {
		r.logger.Debugw("Failed to fetch banner",
			"bannerID", id,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return nil, err
	}
	r.logger.Debugw("Successfully fetched banner",
		"bannerID", id,
		"timeTaken", time.Since(startTime).String(),
	)

	return banner, nil
}

// func (r *BannerRepository) GetRandomBanner(id int) (*entity.Banner, error) {
//      row := r.db.QueryRow("SELECT owner_id, title, description, content, status, balance, link FROM banner WHERE id = $1 AND deleted = FALSE;", id)
//      banner := &entity.Banner{}
//      err := row.Scan(&banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Balance, &banner.Link)
//      if err != nil {
//              return nil, err
//      }
//      return banner, nil
// }

func (r *BannerRepository) DeleteBannerByID(owner, id int) error {
	startTime := time.Now()

	r.logger.Debugw("Starting banner deletion",
		"bannerID", id,
		"ownerID", owner,
	)

	// Проверка существования баннера
	var deleted bool
	err := r.db.QueryRow(
		"SELECT deleted FROM banner WHERE id = $1 AND owner_id = $2",
		id, owner,
	).Scan(&deleted)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugw("Banner not found for deletion",
				"bannerID", id,
				"ownerID", owner,
				"timeTaken", time.Since(startTime).String(),
			)
			return errors.New("banner not found")
		}

		r.logger.Debugw("Failed to check banner existence",
			"bannerID", id,
			"ownerID", owner,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	if deleted {
		r.logger.Debugw("Attempt to delete already deleted banner",
			"bannerID", id,
			"ownerID", owner,
			"timeTaken", time.Since(startTime).String(),
		)
		return errors.New("banner not found")
	}
	// Выполнение удаления
	_, err = r.db.Exec("UPDATE banner SET deleted = TRUE WHERE id = $1", id)
	if err != nil {
		r.logger.Debugw("Failed to delete banner",
			"bannerID", id,
			"ownerID", owner,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	r.logger.Debugw("Banner successfully deleted",
		"bannerID", id,
		"ownerID", owner,
		"timeTaken", time.Since(startTime).String(),
	)
	return nil
}

func (r *BannerRepository) CloseConnection() error {
	return r.db.Close()
}
