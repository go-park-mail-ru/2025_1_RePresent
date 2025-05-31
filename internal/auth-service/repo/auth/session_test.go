package repo

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	authEntity "retarget/internal/auth-service/entity/auth"
)

func setupMiniredis(t *testing.T, ttl time.Duration) (*SessionRepository, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis start failed: %v", err)
	}
	repo := NewSessionRepository(s.Addr(), "", 0, ttl)
	return repo, s
}

func TestGenerateSessionID(t *testing.T) {
	repo, _ := setupMiniredis(t, time.Second)
	id, err := repo.generateSessionID()
	if err != nil || len(id) == 0 {
		t.Fatalf("expected non-empty id, got '%s', err %v", id, err)
	}
}

func TestAddAndGetSession_Success(t *testing.T) {
	repo, _ := setupMiniredis(t, time.Second)
	sess, err := repo.AddSession(42, 2)
	if err != nil {
		t.Fatalf("AddSession error: %v", err)
	}
	got, err2 := repo.GetSession(sess.ID)
	if err2 != nil {
		t.Fatalf("GetSession error: %v", err2)
	}
	if got.UserID != 42 || got.Role != 2 || got.ID != sess.ID {
		t.Errorf("unexpected session %+v", got)
	}
}

func TestGetSession_NotFound(t *testing.T) {
	repo, _ := setupMiniredis(t, time.Second)
	_, err := repo.GetSession("no-such-id")
	if authEntity.ErrSessionNotFound.Error() != err.Error() {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestSessionExpiration(t *testing.T) {
	repo, _ := setupMiniredis(t, 50*time.Millisecond)
	sess, err := repo.AddSession(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	_, err2 := repo.GetSession(sess.ID)
	if err2 != authEntity.ErrSessionNotFound {
		t.Errorf("expected expired, got %v", err2)
	}
}

func TestDelSession(t *testing.T) {
	repo, _ := setupMiniredis(t, time.Second)
	sess, _ := repo.AddSession(7, 1)
	if err := repo.DelSession(sess.ID); err != nil {
		t.Fatalf("DelSession error: %v", err)
	}
	_, err2 := repo.GetSession(sess.ID)
	if err2 != authEntity.ErrSessionNotFound {
		t.Errorf("expected deleted session, got %v", err2)
	}
}

func TestCloseConnection(t *testing.T) {
	repo, s := setupMiniredis(t, time.Second)
	s.Close()
	if err := repo.CloseConnection(); err != nil {
		t.Errorf("CloseConnection error: %v", err)
	}
}
