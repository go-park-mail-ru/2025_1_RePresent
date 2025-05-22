package mail

import (
	"errors"
	entityMail "retarget/internal/mail-service/entity/mail"
	repoMail "retarget/internal/mail-service/repo/mail"
)

type MailUsecaseInterface interface {
	SendCodeMail(operation int, to, code string) error
}

type MailUsecase struct {
	mailRepository *repoMail.MailRepository
}

func NewMailUsecase(mailRepo *repoMail.MailRepository) *MailUsecase {
	return &MailUsecase{mailRepository: mailRepo}
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
	case entityMail.TOPUP_BALANCE:
		subject = "Пополнение баланса в ReTarget"
		body, err = entityMail.GetEmailBody(entityMail.TOPUP_BALANCE, code)
	case entityMail.LOW_BALANCE:
		subject = "Уведомление о низком балансе ReTarget"
		body, err = entityMail.GetEmailBody(entityMail.LOW_BALANCE, code)
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
