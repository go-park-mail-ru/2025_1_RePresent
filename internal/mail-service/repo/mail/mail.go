package mail

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"
)

var (
	ErrInvalidEmail   = errors.New("invalid email address")
	ErrNoMXRecords    = errors.New("no MX records found")
	ErrSMTPConnect    = errors.New("failed to connect to SMTP server")
	ErrEmailNotExists = errors.New("email address does not exist")
	ErrSMTPCommand    = errors.New("SMTP command failed")
)

type MailRepositoryInterface interface {
	Send(to, msg string) error
	VerifyEmail(email string) (bool, error)

	OpenConnection(smtpServer, smtpPort, username, password string) (*smtp.Client, error)
	CloseConnection() error
}

type MailRepository struct {
	smtpClient  *smtp.Client
	from_sender string
}

func NewMailRepository(smtpServer, smtpPort, username, password, from_sender string) *MailRepository {
	mailRepo := &MailRepository{from_sender: from_sender}
	smtpClient, err := mailRepo.OpenConnection(smtpServer, smtpPort, username, password)
	if err != nil {
		log.Fatal(err)
	}
	mailRepo.smtpClient = smtpClient
	return mailRepo
}

func (r *MailRepository) Send(to, msg string) error {
	isExist, err := r.VerifyEmail(to)
	if !isExist || err != nil {
		return err
	}
	if err := r.smtpClient.Mail(r.from_sender); err != nil {
		return err
	}
	if err := r.smtpClient.Rcpt(to); err != nil {
		return err
	}
	w, err := r.smtpClient.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *MailRepository) VerifyEmail(email string) (bool, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, ErrInvalidEmail
	}
	domain := parts[1]

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false, ErrNoMXRecords
	}

	host := mxRecords[0].Host

	conn, err := net.DialTimeout("tcp", host+":25", 10*time.Second)
	if err != nil {
		return false, ErrSMTPConnect
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false, ErrSMTPCommand
	}
	defer client.Quit()

	if err := client.Mail(r.from_sender); err != nil {
		return false, err
	}
	if err := client.Rcpt(email); err != nil {
		return false, ErrEmailNotExists
	}

	return true, nil
}

func (r *MailRepository) OpenConnection(smtpServer, smtpPort, username, password string) (*smtp.Client, error) {
	var smtpClient *smtp.Client

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	email_connect, err := tls.Dial("tcp", smtpServer+":"+smtpPort, tlsConfig)
	if err != nil {
		return nil, err
	}

	smtpClient, err = smtp.NewClient(email_connect, smtpServer)
	if err != nil {
		return nil, err
	}

	err = smtpClient.Auth(smtp.PlainAuth("", username, password, smtpServer))
	if err != nil {
		return nil, err
	}
	return smtpClient, nil
}

func (r *MailRepository) CloseConnection() error {
	return r.smtpClient.Quit()
}
