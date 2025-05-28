package slot

import (
	"context"
	"errors"
	"log"
	"retarget/internal/adv-service/entity/slot"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

var (
	ErrSlotNotFound = errors.New("slot not found")
)

type SlotRepositoryInterface interface {
	CreateSlot(ctx context.Context, userID int, s slot.Slot) (slot.Slot, error)
	UpdateSlot(ctx context.Context, userID int, s slot.Slot) error
	DeleteSlot(ctx context.Context, userID int, link string, created_at time.Time) error
	GetSlotsByUser(ctx context.Context, userID int) ([]slot.Slot, error)
	GetCurrentFormats(ctx context.Context) ([]slot.Format, error)
	GetUserByLink(ctx context.Context, link string) (int, time.Time, error)
	HealthCheck(ctx context.Context) error
	GetSlotInfoByLink(ctx context.Context, link string) (slot.Slot, error)
}

type SlotSession interface {
	Query(stmt string, values ...interface{}) *gocql.Query
	NewBatch(batchType gocql.BatchType) *gocql.Batch
	ExecuteBatch(batch *gocql.Batch) error
	Close()
}

type SlotRepository struct {
	session *gocql.Session
}

func NewSlotRepository(host string, port int, keyspace, username, password string) *SlotRepository {
	cluster := gocql.NewCluster(strings.Split(host, ",")...)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 3 * time.Second
	cluster.ConnectTimeout = 3 * time.Second
	cluster.ReconnectInterval = 100 * time.Millisecond
	cluster.NumConns = 5
	cluster.SocketKeepalive = 15 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	cluster.ReconnectionPolicy = &gocql.ExponentialReconnectionPolicy{
		MaxInterval: 1 * time.Second,
		MaxRetries:  30,
	}

	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(
		gocql.RoundRobinHostPolicy(),
	)

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to create Scylla session: %v", err)
	}

	if err := session.Query("SELECT now() FROM system.local").Exec(); err != nil {
		log.Fatalf("Failed to connect to Scylla: %v", err)
	}

	return &SlotRepository{session: session}
}

func (r *SlotRepository) CreateSlot(ctx context.Context, userID int, s slot.Slot) (slot.Slot, error) {
	batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(
		`INSERT INTO slots (link, user_id, slot_name, format_code, min_price, is_active, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		s.Link, userID, s.SlotName, s.FormatCode, s.MinPrice, s.IsActive, s.CreatedAt,
	)
	batch.Query(
		`INSERT INTO user_links (user_id, created_at, link) VALUES (?, ?, ?)`,
		userID, s.CreatedAt, s.Link,
	)

	if err := r.session.ExecuteBatch(batch); err != nil {
		return slot.Slot{}, err
	}
	return s, nil
}

func (r *SlotRepository) UpdateSlot(ctx context.Context, userID int, s slot.Slot) error {
	applied := false
	var current_user int

	err := r.session.Query(
		`UPDATE slots SET 
			slot_name = ?, 
			format_code = ?, 
			min_price = ?, 
			is_active = ?
		WHERE link = ? IF user_id = ?`,
		s.SlotName, s.FormatCode, s.MinPrice, s.IsActive, s.Link, userID,
	).WithContext(ctx).Scan(&applied, &current_user)

	if err != nil {
		return err
	}
	if !applied {
		return ErrSlotNotFound
	}
	return nil
}

func (r *SlotRepository) DeleteSlot(ctx context.Context, userID int, link string, createdAt time.Time) error {
	batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)

	batch.Query(
		`DELETE FROM slots WHERE link = ?`,
		link,
	)

	batch.Query(
		`DELETE FROM user_links WHERE user_id = ? AND created_at = ? AND link = ?`,
		userID, createdAt, link,
	)

	return r.session.ExecuteBatch(batch)
}

func (r *SlotRepository) GetSlotsByUser(ctx context.Context, userID int) ([]slot.Slot, error) {
	iter := r.session.Query(
		`SELECT link FROM user_links WHERE user_id = ?`,
		userID,
	).WithContext(ctx).Iter()
	defer iter.Close()

	var links []string
	var link string
	for iter.Scan(&link) {
		links = append(links, link)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return []slot.Slot{}, nil
	}

	iter = r.session.Query(
		`SELECT link, slot_name, format_code, min_price, is_active, created_at 
		FROM slots WHERE link IN ?`,
		links,
	).WithContext(ctx).Iter()
	defer iter.Close()

	var slots []slot.Slot
	var s slot.Slot
	for iter.Scan(&s.Link, &s.SlotName, &s.FormatCode, &s.MinPrice, &s.IsActive, &s.CreatedAt) {
		slots = append(slots, s)
	}

	return slots, iter.Close()
}

func (r *SlotRepository) GetCurrentFormats(ctx context.Context) ([]slot.Format, error) {
	iter := r.session.Query(
		`SELECT code, height, width, description FROM formats`,
	).WithContext(ctx).Iter()
	defer iter.Close()

	var formats []slot.Format
	var f slot.Format
	for iter.Scan(&f.Code, &f.Height, &f.Width, &f.Description) {
		formats = append(formats, f)
	}

	return formats, iter.Close()
}

func (r *SlotRepository) GetUserByLink(ctx context.Context, link string) (int, time.Time, error) {
	var userID int
	var createdAt time.Time
	err := r.session.Query(
		`SELECT user_id, created_at FROM slots WHERE link = ? LIMIT 1`,
		link,
	).WithContext(ctx).Scan(&userID, &createdAt)
	if err != nil {
		if err == gocql.ErrNotFound {
			return -1, time.Time{}, ErrSlotNotFound
		}
		return -1, time.Time{}, err
	}
	return userID, createdAt, nil
}

func (r *SlotRepository) GetSlotInfoByLink(ctx context.Context, link string) (slot.Slot, error) {
	var s slot.Slot

	err := r.session.Query(
		`SELECT link, slot_name, format_code, min_price, is_active, created_at 
		FROM slots WHERE link = ? LIMIT 1`,
		link,
	).WithContext(ctx).Scan(
		&s.Link,
		&s.SlotName,
		&s.FormatCode,
		&s.MinPrice,
		&s.IsActive,
		&s.CreatedAt,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			return slot.Slot{}, ErrSlotNotFound
		}
		return slot.Slot{}, err
	}

	return s, nil
}

func (r *SlotRepository) HealthCheck(ctx context.Context) error {
	return r.session.Query("SELECT now() FROM system.local").WithContext(ctx).Exec()
}

func (r *SlotRepository) Close() {
	if r.session != nil {
		r.session.Close()
	}
}
