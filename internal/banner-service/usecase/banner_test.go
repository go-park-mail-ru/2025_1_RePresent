package usecase

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"

	model "retarget/internal/banner-service/easyjsonModels"
	"retarget/internal/banner-service/entity"
	"retarget/internal/banner-service/repo"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newUC(t *testing.T) (*BannerUsecase, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	br := repo.NewBannerRepository("postgres://u:p@/db?sslmode=disable", zap.NewNop().Sugar())
	br.Db = db
	uc := NewBannerUsecase(br)
	uc.rng = rand.New(rand.NewSource(1))
	return uc, mock, func() { db.Close() }
}

func TestGetBannersByUserID_Usecase(t *testing.T) {
	uc, mock, close := newUC(t)
	defer close()

	mock.ExpectQuery("SELECT id, owner_id").
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "title", "description", "content", "status", "link", "max_price"}).
			AddRow(5, 42, "T", "D", "C", 1, "L", "2.0"))
	out, err := uc.GetBannersByUserID(42, "rid")
	assert.NoError(t, err)
	assert.Len(t, out, 1)

	mock.ExpectQuery("SELECT").WithArgs(99).WillReturnError(errors.New("dbFail"))
	_, err = uc.GetBannersByUserID(99, "rid")
	assert.EqualError(t, err, "dbFail")
}

func TestGetBannerByID_Usecase(t *testing.T) {
	uc, mock, close := newUC(t)
	defer close()

	mock.ExpectQuery("SELECT owner_id").
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"owner_id", "title", "description", "content", "balance", "link", "status", "max_price"}).
			AddRow(3, "T", "D", "C", 0, "L", 1, "1.1"))
	b, err := uc.GetBannerByID(3, 7, "rid")
	assert.NoError(t, err)
	assert.Equal(t, 3, b.OwnerID)

	mock.ExpectQuery("SELECT owner_id").WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"owner_id", "title", "description", "content", "balance", "link", "status", "max_price"}).
			AddRow(8, "T", "D", "C", 0, "L", 1, "1.1"))
	_, err = uc.GetBannerByID(3, 7, "rid")
	assert.EqualError(t, err, "banner not found")

	mock.ExpectQuery("SELECT owner_id").WithArgs(7).
		WillReturnError(errors.New("notFound"))
	_, err = uc.GetBannerByID(3, 7, "rid")
	assert.EqualError(t, err, "notFound")
}

func TestGetRandomBannerForIFrame_Usecase(t *testing.T) {
	uc, mock, close := newUC(t)
	defer close()

	mock.ExpectQuery("SELECT id, owner_id").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "title", "description", "content", "status", "link", "max_price"}).
			AddRow(1, 10, "", "", "", 0, "", "0").
			AddRow(2, 10, "", "", "", 0, "", "0"))
	b, _ := uc.GetRandomBannerForIFrame(10, "rid")
	assert.Equal(t, 2, b.ID)

	mock.ExpectQuery("SELECT").WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "title", "description", "content", "status", "link", "max_price"}))
	b2, _ := uc.GetRandomBannerForIFrame(11, "rid")
	assert.Equal(t, entity.DefaultBanner.ID, b2.ID)
}

func TestGetRandomBannerForADV_Usecase(t *testing.T) {
	uc, mock, close := newUC(t)
	defer close()

	mock.ExpectQuery("SELECT b.id").
		WithArgs("0.0").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "description", "link", "owner_id"}).
			AddRow(7, "", "", "", "", 1))
	adv, _ := uc.GetRandomBannerForADV(1, "rid", &entity.DefaultBanner.MaxPrice)
	assert.Equal(t, -1, adv.ID)

	mock.ExpectQuery("SELECT b.id").
		WithArgs("0.0").
		WillReturnError(sql.ErrNoRows)
	adv2, _ := uc.GetRandomBannerForADV(1, "rid", &entity.DefaultBanner.MaxPrice)
	assert.Equal(t, entity.DefaultBanner.ID, adv2.ID)
}

func TestCreateAndUpdateAndDeleteBanner_Usecase(t *testing.T) {
	uc, mock, close := newUC(t)
	defer close()

	mock.ExpectPrepare("INSERT INTO banner").
		ExpectQuery().
		WithArgs(9, "", "", "", 0, 0, "", entity.DefaultBanner.MaxPrice).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
	err := uc.CreateBanner(9, model.Banner{OwnerID: 9}, "rid")
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT owner_id").WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"owner_id", "title", "description", "content", "balance", "link", "status", "max_price"}).
			AddRow(5, "", "", "", 0, "", 0, "0"))
	mock.ExpectPrepare("UPDATE banner").
		ExpectExec().
		WithArgs("", "", "", "", 0, entity.DefaultBanner.MaxPrice, 5).
		WillReturnError(errors.New("uErr"))
	err = uc.UpdateBanner(5, model.Banner{ID: 5, OwnerID: 5}, "rid")
	assert.EqualError(t, err, "uErr")

	mock.ExpectQuery("SELECT deleted FROM banner").WithArgs(5, 5).
		WillReturnRows(sqlmock.NewRows([]string{"deleted"}).AddRow(false))
	mock.ExpectExec("UPDATE banner SET deleted = TRUE").
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = uc.DeleteBannerByID(5, 5, "rid")
	assert.NoError(t, err)
}
