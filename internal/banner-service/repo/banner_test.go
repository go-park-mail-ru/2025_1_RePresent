package repo

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	model "retarget/internal/banner-service/easyjsonModels"
	"retarget/internal/banner-service/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newMockDB(t *testing.T) (*BannerRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	logger := zap.NewNop().Sugar()
	repo := &BannerRepository{Db: db, logger: logger}
	return repo, mock, func() { db.Close() }
}

func TestGetBannersByUserId_SuccessAndErrors(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	cols := []string{"id", "owner_id", "title", "description", "content", "status", "link", "max_price"}
	rows := sqlmock.NewRows(cols).
		AddRow(1, 2, "T", "D", "C", 1, "L", "0.5")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, owner_id")).
		WithArgs(2).
		WillReturnRows(rows)

	out, err := r.GetBannersByUserId(2, "rid")
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, 1, out[0].ID)
	assert.Equal(t, 2, out[0].OwnerID)

	mock.ExpectQuery("SELECT").WithArgs(3).
		WillReturnError(errors.New("qErr"))
	_, err = r.GetBannersByUserId(3, "rid")
	assert.EqualError(t, err, "qErr")

	bad := sqlmock.NewRows([]string{"id"}).
		AddRow(5)
	mock.ExpectQuery("SELECT").WithArgs(4).WillReturnRows(bad)
	_, err = r.GetBannersByUserId(4, "rid")
	assert.Error(t, err)
}

func TestGetMaxPriceBanner(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	mock.ExpectQuery("SELECT b.id").
		WithArgs("1.23").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "description", "link", "owner_id"}).
			AddRow(9, "T", "C", "D", "L", 7))
	out := r.GetMaxPriceBanner(&entity.DefaultBanner.MaxPrice)
	assert.Equal(t, -1, out.ID)

	mock.ExpectQuery("SELECT b.id").
		WithArgs("1.23").
		WillReturnError(sql.ErrNoRows)
	out2 := r.GetMaxPriceBanner(&entity.DefaultBanner.MaxPrice)
	assert.Equal(t, entity.DefaultBanner.ID, out2.ID)
}

func TestCreateNewBanner_PrepareAndExec(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	b := model.Banner{OwnerID: 4, Title: "T", Description: "D", Content: "C", Status: 1, Link: "L", MaxPrice: entity.DefaultBanner.MaxPrice}

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO banner")).
		WillReturnError(errors.New("pErr"))
	err := r.CreateNewBanner(b, "rid")
	assert.EqualError(t, err, "pErr")

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO banner")).
		ExpectQuery().
		WithArgs(4, "T", "D", "C", 1, 0, "L", entity.DefaultBanner.MaxPrice).
		WillReturnError(errors.New("qErr"))
	err = r.CreateNewBanner(b, "rid")
	assert.EqualError(t, err, "qErr")

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO banner")).
		ExpectQuery().
		WithArgs(4, "T", "D", "C", 1, 0, "L", entity.DefaultBanner.MaxPrice).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
	err = r.CreateNewBanner(b, "rid")
	assert.NoError(t, err)
}

func TestUpdateBanner_PrepareAndExec(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	b := model.Banner{ID: 5, Title: "T2", Description: "D2", Content: "C2", Status: 2, Link: "L2", MaxPrice: entity.DefaultBanner.MaxPrice}

	mock.ExpectPrepare("UPDATE banner").
		WillReturnError(errors.New("pErr"))
	err := r.UpdateBanner(b, "rid")
	assert.EqualError(t, err, "pErr")

	mock.ExpectPrepare("UPDATE banner").
		ExpectExec().
		WithArgs("T2", "D2", "C2", "L2", 2, entity.DefaultBanner.MaxPrice, 5).
		WillReturnError(errors.New("eErr"))
	err = r.UpdateBanner(b, "rid")
	assert.EqualError(t, err, "eErr")

	mock.ExpectPrepare("UPDATE banner").
		ExpectExec().
		WithArgs("T2", "D2", "C2", "L2", 2, entity.DefaultBanner.MaxPrice, 5).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = r.UpdateBanner(b, "rid")
	assert.NoError(t, err)
}

func TestGetBannerByID_SuccessAndError(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	mock.ExpectQuery("SELECT owner_id").
		WithArgs(10).
		WillReturnError(sql.ErrNoRows)
	_, err := r.GetBannerByID(10, "rid")
	assert.Error(t, err)

	cols := []string{"owner_id", "title", "description", "content", "balance", "link", "status", "max_price"}
	mock.ExpectQuery("SELECT owner_id").
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(3, "T", "D", "C", 0, "L", 1, "2.2"))
	b, err := r.GetBannerByID(11, "rid")
	assert.NoError(t, err)
	assert.Equal(t, 3, b.OwnerID)
}

func TestDeleteBannerByID_AllBranches(t *testing.T) {
	r, mock, close := newMockDB(t)
	defer close()

	mock.ExpectQuery("SELECT deleted FROM banner").
		WithArgs(8, 7). // изменено: id=8, owner=7
		WillReturnError(sql.ErrNoRows)
	err := r.DeleteBannerByID(7, 8, "rid")
	assert.EqualError(t, err, "banner not found")

	mock.ExpectQuery("SELECT deleted FROM banner").
		WithArgs(8, 7). // изменено: id=8, owner=7
		WillReturnRows(sqlmock.NewRows([]string{"deleted"}).AddRow(true))
	err = r.DeleteBannerByID(7, 8, "rid")
	assert.EqualError(t, err, "banner not found")

	mock.ExpectQuery("SELECT deleted FROM banner").
		WithArgs(8, 7). // изменено: id=8, owner=7
		WillReturnRows(sqlmock.NewRows([]string{"deleted"}).AddRow(false))
	mock.ExpectExec("UPDATE banner SET deleted = TRUE").
		WithArgs(8).
		WillReturnError(errors.New("uErr"))
	err = r.DeleteBannerByID(7, 8, "rid")
	assert.EqualError(t, err, "uErr")

	mock.ExpectQuery("SELECT deleted FROM banner").
		WithArgs(8, 7). // изменено: id=8, owner=7
		WillReturnRows(sqlmock.NewRows([]string{"deleted"}).AddRow(false))
	mock.ExpectExec("UPDATE banner SET deleted = TRUE").
		WithArgs(8).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = r.DeleteBannerByID(7, 8, "rid")
	assert.NoError(t, err)
}
