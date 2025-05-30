package repo

import (
	"database/sql"
	"errors"
	"log"
	model "retarget/internal/banner-service/easyjsonModels"
	"retarget/internal/banner-service/entity"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	decimal "retarget/pkg/entity"
	"time"
)

type BannerRepositoryInterface interface {
	GetBannerByUserId(id int) ([]*model.Banner, error)
	CreateNewBanner(banner model.Banner)
	UpdateBanner(banner model.Banner)
	GetBannerByID(id int) (*model.Banner, error)
	DeleteBannerByID(owner, id int) error
	GetSuitableBanners(floor *decimal.Decimal) ([]int64, error)
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

func (r *BannerRepository) GetSuitableBanners(floor *decimal.Decimal) ([]int64, error) {
	r.logger.Debugw("Executing SQL query GetSuitableBanners", "floor", floor.String())

	query := `
        SELECT b.id
        FROM banner b
        JOIN auth_user u ON b.owner_id = u.id
        WHERE b.status = 1
          AND u.balance > 0
          AND b.max_price >= $1
          AND u.balance >= b.max_price
		  AND NOT b.deleted
    `

	rows, err := r.db.Query(query, floor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []int64{-1}, nil
	}

	return ids, nil
}

func (r *BannerRepository) GetBannersByUserId(id int, requestID string) ([]model.Banner, error) {
	query := "SELECT id, owner_id, title, description, content, status, link, max_price FROM banner WHERE owner_id = $1 AND NOT b.deleted;"
	r.logger.Debugw("Executing SQL query GetProfileByID", "request_id", requestID, "query", query, "userID", id)
	startTime := time.Now()
	rows, err := r.db.Query(query, id)
	duration := time.Since(startTime)
	if err != nil {
		r.logger.Debugw("SQL Error", "request_id", requestID, "userID", id, "duration", duration, "error", err)
		return nil, err
	}
	defer rows.Close()
	banners := []model.Banner{}

	for rows.Next() {
		banner := model.Banner{}
		err := rows.Scan(&banner.ID, &banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Link, &banner.MaxPrice)
		if err != nil {
			r.logger.Debugw("SQL Error", "request_id", requestID, "userID", id, "duration", duration, "error", err)
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		r.logger.Debugw("SQL Error", "request_id", requestID, "userID", id, "duration", duration, "error", err)
		return nil, err
	}
	r.logger.Debugw("SQL  query executed successfully", "request_id", requestID, "userID", id, "duration", duration, "error", err)
	return banners, nil
}

func (r *BannerRepository) GetMaxPriceBanner(floor *decimal.Decimal) *model.Banner {
	r.logger.Debugw("Executing SQL query GetMaxPriceBanner", "floor", floor.String())

	query := `
        SELECT b.id, b.title, b.content, b.description, b.link, b.owner_id
        FROM banner b
        JOIN auth_user u ON b.owner_id = u.id
        WHERE b.status = 1
          AND u.balance > 0
          AND b.max_price = (
              SELECT MAX(b2.max_price)
              FROM banner b2
              JOIN auth_user u2 ON b2.owner_id = u2.id
              WHERE b2.status = 1 AND u2.balance > 0 AND b2.max_price > $1
          )
		  AND NOT b.deleted
        ORDER BY RANDOM()
        LIMIT 1;
    `

	var banner model.Banner
	startTime := time.Now()

	err := r.db.QueryRow(query, floor).Scan(
		&banner.ID,
		&banner.Title,
		&banner.Content,
		&banner.Description,
		&banner.Link,
		&banner.OwnerID,
	)

	r.logger.Debugw("SQL query completed", "duration", time.Since(startTime))

	if err != nil {
		return &entity.DefaultBanner
	}

	return &banner
}

func (r *BannerRepository) CreateNewBanner(banner model.Banner, requestID string) error {
	r.logger.Debugw("Executing SQL query СreateNewBanner",
		"request_id", requestID,
		"ownerID", banner.OwnerID,
		// "title", banner.Title,
		// "description", banner.Description,
		// "status", banner.Status,
		"link", banner.Link,
	)
	stmt, err := r.db.Prepare("INSERT INTO banner (owner_id, title, description, content, status, balance, link, max_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;")
	startTime := time.Now()

	if err != nil {
		r.logger.Debugw("Error preparing SQL statement",
			"request_id", requestID, "error", err,
		)
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(banner.OwnerID, banner.Title, banner.Description, banner.Content, banner.Status, 0, banner.Link, banner.MaxPrice).Scan(&id)
	if err != nil {
		r.logger.Debugw("Error executing query to create new banner", "request_id", requestID, "error", err)
		return err
	}
	duration := time.Since(startTime)
	r.logger.Debugw("Successfully created new banner", "request_id", requestID, "bannerID", id, "duration", duration)

	return nil

}

func (r *BannerRepository) UpdateBanner(banner model.Banner, requestID string) error {
	startTime := time.Now()
	query := "UPDATE banner SET title = $1, description = $2, content = $3, link = $4, status = $5, max_price = $6 WHERE id = $7"
	r.logger.Debugw("Starting banner update",
		"request_id", requestID,
		"bannerID", banner.ID,
		"query", query,
	)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		r.logger.Debugw("Failed to prepare SQL statement for banner update",
			"request_id", requestID,
			"error", err,
			"bannerID", banner.ID,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(banner.Title, banner.Description, banner.Content, banner.Link, banner.Status, banner.MaxPrice, banner.ID)
	if err != nil {
		r.logger.Debugw("Failed to execute banner update",
			"request_id", requestID,
			"error", err,
			"bannerID", banner.ID,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	r.logger.Debugw("Banner updated successfully",
		"request_id", requestID,
		"bannerID", banner.ID,
		"timeTaken", time.Since(startTime).String(),
	)
	return nil
}

func (r *BannerRepository) GetBannerByID(id int, requestID string) (*model.Banner, error) {
	startTime := time.Now()
	query := `
		SELECT owner_id, title, description, content, balance, link, status, max_price
		FROM banner
		WHERE id = $1 AND deleted = FALSE;
		`
	r.logger.Debugw("Starting to fetch banner by ID",
		"request_id", requestID,
		"bannerID", id, "query", query,
	)

	row := r.db.QueryRow(query, id)

	banner := &model.Banner{}
	err := row.Scan(
		&banner.OwnerID,
		&banner.Title,
		&banner.Description,
		&banner.Content,
		&banner.Balance,
		&banner.Link,
		&banner.Status,
		&banner.MaxPrice,
	)
	if err != nil {
		r.logger.Debugw("Failed to fetch banner",
			"request_id", requestID,
			"bannerID", id,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return nil, err
	}
	r.logger.Debugw("Successfully fetched banner",
		"request_id", requestID,
		"bannerID", id,
		"timeTaken", time.Since(startTime).String(),
	)

	return banner, nil
}

// func (r *BannerRepository) GetRandomBanner(id int) (*model.Banner, error) {
//      row := r.db.QueryRow("SELECT owner_id, title, description, content, status, balance, link FROM banner WHERE id = $1 AND deleted = FALSE;", id)
//      banner := &model.Banner{}
//      err := row.Scan(&banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Balance, &banner.Link)
//      if err != nil {
//              return nil, err
//      }
//      return banner, nil
// }

func (r *BannerRepository) DeleteBannerByID(owner, id int, requestID string) error {
	startTime := time.Now()

	r.logger.Debugw("Starting banner deletion",
		"request_id", requestID,
		"bannerID", id,
		"ownerID", owner,
	)

	var deleted bool
	err := r.db.QueryRow(
		"SELECT deleted FROM banner WHERE id = $1 AND owner_id = $2",
		id, owner,
	).Scan(&deleted)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugw("Banner not found for deletion",
				"request_id", requestID,
				"bannerID", id,
				"ownerID", owner,
				"timeTaken", time.Since(startTime).String(),
			)
			return errors.New("banner not found")
		}

		r.logger.Debugw("Failed to check banner existence",
			"request_id", requestID,
			"bannerID", id,
			"ownerID", owner,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	if deleted {
		r.logger.Debugw("Attempt to delete already deleted banner",
			"request_id", requestID,
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
			"request_id", requestID,
			"bannerID", id,
			"ownerID", owner,
			"error", err,
			"timeTaken", time.Since(startTime).String(),
		)
		return err
	}
	r.logger.Debugw("Banner successfully deleted",
		"request_id", requestID,
		"bannerID", id,
		"ownerID", owner,
		"timeTaken", time.Since(startTime).String(),
	)
	return nil
}

func (r *BannerRepository) CloseConnection() error {
	return r.db.Close()
}
