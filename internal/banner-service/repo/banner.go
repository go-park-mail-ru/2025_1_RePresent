package repo

import (
	"database/sql"

	"retarget/internal/banner-service/entity"

	_ "github.com/lib/pq"

	"errors"
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
	rows, err := r.db.Query("SELECT id, owner_id, title, description, content, status FROM banner WHERE owner_id = $1 AND deleted = FALSE", id)
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

	stmt, err := r.db.Prepare("INSERT INTO banner (owner, title, description, content, status, balance, link) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(banner.OwnerID, banner.Title, banner.Description, banner.Content, banner.Status, 0, banner.Link).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}

func (r *BannerRepository) GetBannerByID(id int) (*entity.Banner, error) {
	row := r.db.QueryRow("SELECT owner, title, description, content, status, balance, link FROM auth_user WHERE id = $1", id)
	banner := &entity.Banner{}
	err := row.Scan(&banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Balance, &banner.Link)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *BannerRepository) DeleteBannerByID(id, owner int) error {
	var deleted bool
	err := r.db.QueryRow("SELECT deleted FROM banner WHERE id = $1 AND owner = $2", id, owner).Scan(&deleted)

	if err == sql.ErrNoRows {
		return errors.New("Banner not found")
	}
	if err != nil {
		return err
	}
	if deleted {
		return errors.New("Banner not found")
	}

	_, err = r.db.Exec("UPDATE banner SET deleted = TRUE WHERE id = $1", id)
	return err
}
