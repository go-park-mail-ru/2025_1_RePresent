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
			action, 
		) VALUES (?, ?, ?)
	`

	if _, err := u.clickhouse.Exec(addQuery, bannerID, slotLink, action); err != nil {
		return fmt.Errorf("Ошибка при добавлении в базу данных")
	}

	return nil
}
