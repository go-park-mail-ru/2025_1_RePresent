package kafka

import (
	"context"
	"encoding/json"
	"log"
	"retarget/pkg/entity/notice"

	"retarget/internal/mail-service/entity/mail"
	usecaseMail "retarget/internal/mail-service/usecase/mail"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type Consumer struct {
	processor   *goka.Processor
	mailUseCase *usecaseMail.MailUsecase
}

func NewConsumer(brokers []string, group string, topic goka.Stream, mailUseCase *usecaseMail.MailUsecase) *Consumer {
	cb := func(ctx goka.Context, msg interface{}) {
		log.Printf("Raw message: %v", msg)

		var event notice.NoticeEvent
		err := json.Unmarshal([]byte(msg.(string)), &event)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}

		log.Printf("\n[KAFKA] NEW MESSAGE\n"+
			"User ID: %d\n"+
			"Type: %d\n"+
			"-------------------------",
			event.UserID, event.Type)

		email := "froloff1830@gmail.com"
		username := "UserName"
		balance := "100.00"
		href := "https://re-target.ru/profile"
		switch event.Type {
		case notice.LowBalance:
			if err := mailUseCase.SendLowBalanceMail(mail.LOW_BALANCE, email, username, balance, href); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email successfully sent to %s", email)
			}
		case notice.TopUpedBalance:
			if err := mailUseCase.SendTopUpBalanceMail(mail.TOPUP_BALANCE, email, username, balance); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email successfully sent to %s", email)
			}
		default:
			log.Printf("!!! UNDEFINED EVENT IN KAFKA !!!: %v", event.Type)
		}
	}

	input := goka.Input(topic, new(codec.String), cb)

	processor, err := goka.NewProcessor(brokers,
		goka.DefineGroup(goka.Group(group),
			input,
		),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)

	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	return &Consumer{
		processor:   processor,
		mailUseCase: mailUseCase,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	if err := c.processor.Run(ctx); err != nil {
		log.Printf("error running processor: %v", err)
	}
}
