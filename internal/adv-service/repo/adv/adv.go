package adv

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"retarget/internal/adv-service/entity/adv"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gocql/gocql"
)

type AdvRepositoryInterface interface {
	FindUserByLink(link string) (int, error)
	CreateLink(link adv.Link) error
	FindLinksByUser(userID int) ([]adv.Link, error)
	DeleteLink(link string) error
	WriteMetric(bannerID int, slotLink string, action string) error
	GetSlotMetric(slotID, action string, from, to time.Time) (map[string]int, error)
	GetBannerMetric(bannerID int, action string, from, to time.Time) (map[string]int, error)
}

type AdvRepository struct {
	session    *gocql.Session
	clickhouse *sql.DB
}

func NewAdvRepository(host string, port int, keyspace, username, password string, dsn string) *AdvRepository {
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping DB: %w", err)
		return nil
	}

	return &AdvRepository{session: session, clickhouse: db}
}

func (u *AdvRepository) FindUserByLink(link string) (int, error) {
	var userID int
	err := u.session.Query(`SELECT user_id FROM links WHERE link = ?`, link).Scan(&userID)
	if err != nil {
		if err == gocql.ErrNotFound {
			return -1, nil
		}
		return 0, err
	}
	return userID, nil
}

func (u *AdvRepository) CreateLink(link adv.Link) error {
	err := u.session.Query(`INSERT INTO links (link, user_id, height, width) VALUES (?, ?, ?, ?)`,
		link.TextLink, link.UserID, link.Height, link.Width).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *AdvRepository) FindLinksByUser(userID int) ([]adv.Link, error) {
	iter := u.session.Query(`SELECT link, user_id, height, width FROM links WHERE user_id = ?`, userID).Iter()

	var links []adv.Link
	var link adv.Link
	for iter.Scan(&link.TextLink, &link.UserID, &link.Height, &link.Width) {
		links = append(links, link)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return links, nil
}

func (u *AdvRepository) DeleteLink(link string) error {
	err := u.session.Query(`DELETE FROM links WHERE link = ?`, link).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *AdvRepository) WriteMetric(bannerID int, slotLink string, action string) error {
	const addQuery = `
		INSERT INTO actions (
			banner_id, 
			slot_id, 
			actions
		) VALUES (?, ?, ?)
	`

	res, err := u.clickhouse.Exec(addQuery, bannerID, slotLink, action)
	if err != nil {
		log.Printf("ClickHouse insert error: %v", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
	} else {
		log.Printf("Rows affected: %d", rowsAffected)
	}

	return nil
}

func (u *AdvRepository) GetSlotMetric(slotID, action string, from, to time.Time) (map[string]int, error) {
	const query = `
		SELECT toDate(created_at) as day, count(*) as total
		FROM adv.actions
		WHERE slot_id = ?
		AND actions = ?
		AND created_at BETWEEN ? AND ?
		GROUP BY day
		ORDER BY day
	`

	rows, err := u.clickhouse.Query(query, slotID, action, from, to)
	if err != nil {
		return nil, fmt.Errorf("error when reading from the database")
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var (
			date  time.Time
			count int
		)

		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("error when reading from the database")
		}

		result[date.Format("2006-01-02")] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when reading from the database")
	}

	return result, nil
}

func (u *AdvRepository) GetBannerMetric(bannerID int, action string, from, to time.Time) (map[string]int, error) {
	const query = `
		SELECT toDate(created_at) as day, count(*) as total
		FROM adv.actions
		WHERE banner_id = ?
		AND actions = ?
		AND created_at BETWEEN ? AND ?
		GROUP BY day
		ORDER BY day
	`

	rows, err := u.clickhouse.Query(query, bannerID, action, from, to)
	if err != nil {
		return nil, fmt.Errorf("error when reading from the database")
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var (
			date  time.Time
			count int
		)

		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("error when reading from the database")
		}

		result[date.Format("2006-01-02")] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when reading from the database")
	}

	return result, nil
}
