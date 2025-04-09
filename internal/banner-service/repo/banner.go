package repo

import (
	"database/sql"

	"log"
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

func NewBannerRepository(endPoint string) *BannerRepository {
	bannerRepo := &BannerRepository{}
	db, err := sql.Open("postgres", endPoint)
	if err != nil {
		log.Fatal(err)
	}
	bannerRepo.db = db
	return bannerRepo
}

func (r *BannerRepository) GetBannersByUserId(id int) ([]*entity.Banner, error) {
	rows, err := r.db.Query("SELECT id, owner_id, title, description, content, status, link FROM banner WHERE owner_id = $1 AND deleted = FALSE;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var banners []*entity.Banner

	for rows.Next() {
		banner := &entity.Banner{}
		err := rows.Scan(&banner.ID, &banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Status, &banner.Link)
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

	stmt, err := r.db.Prepare("INSERT INTO banner (owner_id, title, description, content, status, balance, link) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;")
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

func (r *BannerRepository) UpdateBanner(banner entity.Banner) error {
	stmt, err := r.db.Prepare("UPDATE banner SET title = $1, description = $2, content = $3, link = $4, status = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(banner.Title, banner.Description, banner.Content, banner.Link, banner.Status, banner.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *BannerRepository) GetBannerByID(id int) (*entity.Banner, error) {
	row := r.db.QueryRow(`
    SELECT owner_id, title, description, content, balance, link, status
    FROM banner
    WHERE id = $1 AND deleted = FALSE;
	`, id)
	banner := &entity.Banner{}
	err := row.Scan(&banner.OwnerID, &banner.Title, &banner.Description, &banner.Content, &banner.Balance, &banner.Link, &banner.Status)
	if err != nil {
		return nil, err
	}
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
	var deleted bool
	err := r.db.QueryRow("SELECT deleted FROM banner WHERE id = $1 AND owner_id = $2", id, owner).Scan(&deleted)

	if err == sql.ErrNoRows {
		return errors.New("banner not found")
	}
	if err != nil {
		return err
	}
	if deleted {
		return errors.New("banner not found")
	}

	_, err = r.db.Exec("UPDATE banner SET deleted = TRUE WHERE id = $1", id)
	return err
}

func (r *BannerRepository) CloseConnection() error {
	return r.db.Close()
}
