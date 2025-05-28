package adv

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockAdvRepo(t *testing.T) (*AdvRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return &AdvRepository{clickhouse: db, session: nil}, mock
}

func TestWriteMetric_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectExec("INSERT INTO actions").
		WithArgs(1, "slot1", "click", "100").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.WriteMetric(1, "slot1", "click", "100"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestWriteMetric_ExecError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectExec("INSERT INTO actions").
		WithArgs(1, "s", "a", "p").
		WillReturnError(fmt.Errorf("exec err"))

	if err := repo.WriteMetric(1, "s", "a", "p"); err == nil {
		t.Error("expected error on Exec")
	}
}

func TestWriteMetric_RowsAffectedError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	errRes := sqlmock.NewErrorResult(fmt.Errorf("ra err"))
	mock.ExpectExec("INSERT INTO actions").
		WithArgs(2, "l", "a2", "pr").
		WillReturnResult(errRes)

	if err := repo.WriteMetric(2, "l", "a2", "pr"); err != nil {
		t.Errorf("expected no error despite RowsAffected failure, got %v", err)
	}
}

func TestGetSlotMetric_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total"}).
		AddRow(from, 5).
		AddRow(to, 7)
	mock.ExpectQuery("toDate\\(created_at\\) as day").
		WithArgs("slot1", "shown", from, to).
		WillReturnRows(rows)

	m, err := repo.GetSlotMetric("slot1", "shown", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(m) != 2 || m[from.Format("2006-01-02")] != 5 {
		t.Errorf("bad result %v", m)
	}
}

func TestGetSlotMetric_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("toDate\\(created_at\\)").
		WillReturnError(sql.ErrConnDone)

	if _, err := repo.GetSlotMetric("s", "a", time.Now(), time.Now()); err == nil {
		t.Error("expected error")
	}
}

func TestGetSlotMetric_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total"}).
		AddRow("bad-date", "bad-count")
	mock.ExpectQuery("toDate\\(created_at\\)").
		WithArgs("x", "y", from, to).
		WillReturnRows(rows)

	if _, err := repo.GetSlotMetric("x", "y", from, to); err == nil {
		t.Error("expected scan error for bad data")
	}
}

func TestGetSlotCTR_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "ctr"}).
		AddRow(from, 0.1234).
		AddRow(to, 0.5678)
	mock.ExpectQuery("round\\(clicks / shown").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	res, err := repo.GetSlotCTR("slot1", "ignored", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[from.Format("2006-01-02")] != 0.1234 {
		t.Errorf("got %v", res)
	}
}

func TestGetSlotCTR_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("round\\(clicks / shown").
		WillReturnError(fmt.Errorf("q err"))

	if _, err := repo.GetSlotCTR("s", "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetSlotCTR")
	}
}

func TestGetSlotCTR_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "ctr"}).
		AddRow("bad-date", "bad-ctr")
	mock.ExpectQuery("round\\(clicks / shown").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	if _, err := repo.GetSlotCTR("slot1", "a", from, to); err == nil {
		t.Error("expected scan error for CTR")
	}
}

func TestGetSlotRevenue_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from, to := time.Now(), time.Now().Add(time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total_price"}).
		AddRow(from, 11.11).
		AddRow(to, 22.22)
	mock.ExpectQuery("sum\\(price\\)").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	res, err := repo.GetSlotRevenue("slot1", "ignored", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[to.Format("2006-01-02")] != 22.22 {
		t.Errorf("got %v", res)
	}
}

func TestGetSlotRevenue_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("sum\\(price\\)").
		WillReturnError(fmt.Errorf("rev err"))

	if _, err := repo.GetSlotRevenue("s", "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetSlotRevenue")
	}
}

func TestGetSlotRevenue_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total_price"}).
		AddRow("bad-date", "bad-price")
	mock.ExpectQuery("sum\\(price\\)").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	if _, err := repo.GetSlotRevenue("slot1", "a", from, to); err == nil {
		t.Error("expected scan error for revenue")
	}
}

func TestGetSlotAVGPrice_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "avg_price"}).
		AddRow(from, 3.33).
		AddRow(to, 4.44)
	mock.ExpectQuery("avg\\(price\\)").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	res, err := repo.GetSlotAVGPrice("slot1", "ignored", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[from.Format("2006-01-02")] != 3.33 {
		t.Errorf("got %v", res)
	}
}

func TestGetSlotAVGPrice_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("avg\\(price\\)").
		WillReturnError(fmt.Errorf("avg err"))

	if _, err := repo.GetSlotAVGPrice("s", "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetSlotAVGPrice")
	}
}

func TestGetSlotAVGPrice_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "avg_price"}).
		AddRow("bad-date", "bad-avg")
	mock.ExpectQuery("avg\\(price\\)").
		WithArgs("slot1", from, to).
		WillReturnRows(rows)

	if _, err := repo.GetSlotAVGPrice("slot1", "a", from, to); err == nil {
		t.Error("expected scan error for avg price")
	}
}

func TestGetBannerMetric_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total"}).
		AddRow(from, 9).
		AddRow(to, 8)
	mock.ExpectQuery("WHERE banner_id =").
		WithArgs(42, "click", from, to).
		WillReturnRows(rows)

	res, err := repo.GetBannerMetric(42, "click", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[from.Format("2006-01-02")] != 9 {
		t.Errorf("got %v", res)
	}
}

func TestGetBannerMetric_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("WHERE banner_id =").
		WillReturnError(fmt.Errorf("bm err"))

	if _, err := repo.GetBannerMetric(0, "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetBannerMetric")
	}
}

func TestGetBannerMetric_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total"}).
		AddRow("bad-date", "bad-total")
	mock.ExpectQuery("WHERE banner_id =").
		WithArgs(1, "a", from, to).
		WillReturnRows(rows)

	if _, err := repo.GetBannerMetric(1, "a", from, to); err == nil {
		t.Error("expected scan error for banner metric")
	}
}

func TestGetBannerCTR_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from, to := time.Now(), time.Now().Add(time.Hour)
	rows := sqlmock.NewRows([]string{"day", "ctr"}).
		AddRow(from, 0.5)
	mock.ExpectQuery("round\\(clicks / shown").
		WithArgs(42, from, to).
		WillReturnRows(rows)

	res, err := repo.GetBannerCTR(42, "ignored", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[from.Format("2006-01-02")] != 0.5 {
		t.Errorf("got %v", res)
	}
}

func TestGetBannerCTR_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("round\\(clicks / shown").
		WillReturnError(fmt.Errorf("bc err"))

	if _, err := repo.GetBannerCTR(0, "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetBannerCTR")
	}
}

func TestGetBannerCTR_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "ctr"}).
		AddRow("bad-date", "bad-ctr")
	mock.ExpectQuery("round\\(clicks / shown").
		WithArgs(42, from, to).
		WillReturnRows(rows)

	if _, err := repo.GetBannerCTR(42, "a", from, to); err == nil {
		t.Error("expected scan error for banner CTR")
	}
}

func TestGetBannerExpenses_Success(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from, to := time.Now(), time.Now().Add(time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total_price"}).
		AddRow(to, 7.77)
	mock.ExpectQuery("sum\\(price\\)").
		WithArgs(42, from, to).
		WillReturnRows(rows)

	res, err := repo.GetBannerExpenses(42, "ignored", from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res[to.Format("2006-01-02")] != 7.77 {
		t.Errorf("got %v", res)
	}
}

func TestGetBannerExpenses_QueryError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	mock.ExpectQuery("sum\\(price\\)").
		WillReturnError(fmt.Errorf("be err"))

	if _, err := repo.GetBannerExpenses(0, "a", time.Now(), time.Now()); err == nil {
		t.Error("expected Query error on GetBannerExpenses")
	}
}

func TestGetBannerExpenses_ScanError(t *testing.T) {
	repo, mock := newMockAdvRepo(t)
	from := time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	rows := sqlmock.NewRows([]string{"day", "total_price"}).
		AddRow("bad-date", "bad-exp")
	mock.ExpectQuery("sum\\(price\\)").
		WithArgs(42, from, to).
		WillReturnRows(rows)

	if _, err := repo.GetBannerExpenses(42, "a", from, to); err == nil {
		t.Error("expected scan error for banner expenses")
	}
}
