package mail

import (
	"errors"
	repoAuth "retarget/internal/auth-service/repo/auth"
	entityMail "retarget/internal/mail-service/entity/mail"
	repoMail "retarget/internal/mail-service/repo/mail"
)

type MailUsecaseInterface interface {
	SendCodeMail(operation int, to, code string) error
	GetUserByID(user_id int) (string, string, string, error)
}

type MailUsecase struct {
	mailRepository *repoMail.MailRepository
	userRepository *repoAuth.AuthRepository
}

func NewMailUsecase(mailRepo *repoMail.MailRepository, userRepo *repoAuth.AuthRepository) *MailUsecase {
	return &MailUsecase{mailRepository: mailRepo, userRepository: userRepo}
}

func (m *MailUsecase) SendLowBalanceMail(operation int, to, username, balance, href string) error {
	var subject string
	var body string
	var err error

	switch operation {
	case entityMail.LOW_BALANCE:
		subject = "Уведомление о низком балансе ReTarget"
		body, err = entityMail.GetEmailLowBalanceBody(entityMail.LOW_BALANCE, username, balance, href)
	default:
		return errors.New("undefined operation")
	}

	if err != nil {
		return err
	}

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body

	err = m.mailRepository.Send(to, msg)
	if err != nil {
		return err
	}
	return nil
}

func (m *MailUsecase) SendTopUpBalanceMail(operation int, to, username, amount string) error {
	var subject string
	var body string
	var err error

	switch operation {
	case entityMail.TOPUP_BALANCE:
		subject = "Пополнение баланса ReTarget"
		body, err = entityMail.GetEmailTopUpBody(entityMail.TOPUP_BALANCE, username, amount)
	default:
		return errors.New("undefined operation")
	}

	if err != nil {
		return err
	}

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body

	err = m.mailRepository.Send(to, msg)
	if err != nil {
		return err
	}
	return nil
}

func (m *MailUsecase) SendCodeMail(operation int, to, code string) error {
	var subject string
	var body string
	var err error

	switch operation {
	case entityMail.REGISTER:
		subject = "Подтверждение почты в ReTarget"
		body, err = entityMail.GetEmailBody(entityMail.REGISTER, code)
	case entityMail.RESET_PASSWORD:
		subject = "Сброс пароля в ReTarget"
		body, err = entityMail.GetEmailBody(entityMail.RESET_PASSWORD, code)
	case entityMail.EDIT_PASSWORD:
		subject = "Изменение пароля в ReTarget"
		body, err = entityMail.GetEmailBody(entityMail.EDIT_PASSWORD, code)
	default:
		return errors.New("undefined operation")
	}

	if err != nil {
		return err
	}

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body

	err = m.mailRepository.Send(to, msg)
	if err != nil {
		return err
	}
	return nil
}

func (m *MailUsecase) GetUserByID(user_id int) (string, string, string, error) {
	user, err := m.userRepository.GetUserByID(user_id, "UNIMPIMENTED request_id")
	if err != nil {
		return "froloff1830@gmail.com", "Error if parse Username", "Balance", err
	}
	return user.Email, user.Username, user.Balance.String(), nil
}
