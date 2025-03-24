package repo

import (
	"database/sql"

	"retarget/internal/banner-service/entity"

	_ "github.com/lib/pq"
)

type BannerRepositoryInterface interface {
	GerBannerByUserId(id int) ([]*entity.Banner, error)
}

type BannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) GetBannersByUserId(id int) ([]*entity.Banner, error) {
	rows, err := r.db.Query("SELECT id, owner_id, title, description, content, status FROM banner WHERE owner_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var banners []*entity.Banner

	for rows.Next() {
		banner := &entity.Banner{}
		err := rows.Scan(&banner.ID, &banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *BannerRepository) CreateNewBanner(banner entity.Banner) error {

	stmt, err := r.db.Prepare("INSERT INTO banner (owner, title, description, content, status, balance) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(banner.OwnerID, banner.Title, banner.Description, banner.Content, banner.Status, 0).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}
