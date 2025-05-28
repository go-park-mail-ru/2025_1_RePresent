package slot_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"gopkg.in/inf.v0"

	entity "retarget/internal/adv-service/entity/slot"
	repo "retarget/internal/adv-service/repo/slot"
)

// fakeBatch, fakeQuery, fakeSession для имитации gocql.Session
type fakeBatch struct {
	stmts []string
	args  [][]interface{}
}

func (b *fakeBatch) WithContext(ctx context.Context) repo.Batch { return b }
func (b *fakeBatch) Query(stmt string, values ...interface{}) {
	b.stmts = append(b.stmts, stmt)
	b.args = append(b.args, values)
}

type fakeQuery struct {
	session *fakeSession
	stmt    string
	values  []interface{}
}

func (q *fakeQuery) WithContext(ctx context.Context) repo.Query { return q }
func (q *fakeQuery) Scan(dest ...interface{}) error {
	if q.session.scanErr != nil {
		return q.session.scanErr
	}
	if len(dest) >= 1 {
		*dest[0].(*bool) = q.session.applied
	}
	return nil
}
func (q *fakeQuery) Iter() repo.Iter { return nil }

type fakeSession struct {
	lastBatch *fakeBatch
	applied   bool
	scanErr   error
}

func (s *fakeSession) NewBatch(batchType int) repo.Batch {
	s.lastBatch = &fakeBatch{}
	return s.lastBatch
}
func (s *fakeSession) ExecuteBatch(b repo.Batch) error { return nil }
func (s *fakeSession) Query(stmt string, values ...interface{}) repo.Query {
	return &fakeQuery{session: s, stmt: stmt, values: values}
}
func (s *fakeSession) Close() {}

// newTestRepo создаёт SlotRepository с внедрённым fakeSession
func newTestRepo(sess repo.SlotSession) *repo.SlotRepository {
	r := new(repo.SlotRepository)
	v := reflect.ValueOf(r).Elem().FieldByName("session")
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).
		Elem().Set(reflect.ValueOf(sess))
	return r
}

// TestCreateSlot
func TestCreateSlot(t *testing.T) {
	sess := &fakeSession{}
	r := newTestRepo(sess)

	now := time.Now()
	price := inf.NewDec(100, 0)
	s := entity.Slot{
		Link:       "L",
		SlotName:   "name",
		FormatCode: 1,
		MinPrice:   *price,
		IsActive:   true,
		CreatedAt:  now,
	}

	got, err := r.CreateSlot(context.Background(), 42, s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Link != s.Link ||
		got.SlotName != s.SlotName ||
		got.FormatCode != s.FormatCode ||
		got.IsActive != s.IsActive ||
		got.MinPrice.Cmp(&s.MinPrice) != 0 ||
		!got.CreatedAt.Equal(s.CreatedAt) {
		t.Errorf("got %+v, want %+v", got, s)
	}
	if sess.lastBatch == nil || len(sess.lastBatch.stmts) != 2 {
		t.Errorf("expected 2 queries in batch, got %v", sess.lastBatch)
	}
}

// TestUpdateSlot_Applied
func TestUpdateSlot_Applied(t *testing.T) {
	sess := &fakeSession{applied: true}
	r := newTestRepo(sess)

	s := entity.Slot{Link: "L", SlotName: "n", FormatCode: 1, MinPrice: *inf.NewDec(0, 0), IsActive: false}
	if err := r.UpdateSlot(context.Background(), 1, s); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestUpdateSlot_NotApplied
func TestUpdateSlot_NotApplied(t *testing.T) {
	sess := &fakeSession{applied: false}
	r := newTestRepo(sess)

	s := entity.Slot{Link: "L", SlotName: "n", FormatCode: 1, MinPrice: *inf.NewDec(0, 0), IsActive: false}
	err := r.UpdateSlot(context.Background(), 1, s)
	if !errors.Is(err, repo.ErrSlotNotFound) {
		t.Errorf("expected ErrSlotNotFound, got %v", err)
	}
}

// TestDeleteSlot
func TestDeleteSlot(t *testing.T) {
	sess := &fakeSession{}
	r := newTestRepo(sess)

	now := time.Now()
	if err := r.DeleteSlot(context.Background(), 7, "L", now); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if sess.lastBatch == nil || len(sess.lastBatch.stmts) != 2 {
		t.Errorf("expected 2 delete queries, got %v", sess.lastBatch)
	}
}
