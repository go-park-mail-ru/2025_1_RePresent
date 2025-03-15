package mail

import (
	"crypto/tls"
	"log"
	"net/smtp"
)

type MailRepositoryInterface interface {
	Send(to, msg string) error
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
