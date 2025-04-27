package repo

import (
	"log"

	"retarget/internal/adv-service/entity/adv"

	"github.com/gocql/gocql"
)

type AdvRepositoryInterface interface {
	FindUserByLink(link string) (int, error)
	CreateLink(link adv.Link) error
	FindLinksByUser(userID int) ([]adv.Link, error)
	DeleteLink(link string) error
}

type AdvRepository struct {
	session *gocql.Session
}

func NewAdvRepository(host string, port int, keyspace, username, password string) *AdvRepository {
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
	return &AdvRepository{session: session}
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
