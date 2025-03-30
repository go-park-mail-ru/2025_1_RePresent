package repo

import (
	"log"

	"github.com/gocql/gocql"
)

type AdvRepositoryInterface interface {
	FindUserByLink(secret_link string) (int, error)
	CreateLinkByUser(user_id int, secret_link string) error
	FindLinkByUser(user_id int) (string, error)
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

func (u *AdvRepository) FindUserByLink(secret_link string) (int, error) {
	var user_id int
	err := u.session.Query(`SELECT user_id FROM links WHERE link =?`, secret_link).Scan(&user_id)
	if err != nil {
		if err == gocql.ErrNotFound {
			return -1, gocql.ErrNotFound
		}
		return 0, err
	}
	return user_id, nil
}

func (u *AdvRepository) CreateLinkByUser(user_id int, secret_link string) error {
	err := u.session.Query(`INSERT INTO links (link, user_id) VALUES (?,?)`, secret_link, user_id).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *AdvRepository) FindLinkByUser(user_id int) (string, error) {
	var secret_link string
	err := u.session.Query(`SELECT link FROM links WHERE user_id =?`, user_id).Scan(&secret_link)
	if err != nil {
		if err == gocql.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return secret_link, nil
}
