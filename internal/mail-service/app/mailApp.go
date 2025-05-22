package mailApp

import (
	"context"
	"log"
	"net/http"
	"retarget/configs"
	mailAppHttp "retarget/internal/mail-service/controller/http"
	mailMiddleware "retarget/internal/mail-service/controller/http/middleware"
	mailAppKafka "retarget/internal/mail-service/controller/kafka"
	repoMail "retarget/internal/mail-service/repo/mail"
	usecaseMail "retarget/internal/mail-service/usecase/mail"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	mailRepository := repoMail.NewMailRepository(cfg.Email.SmtpServer, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, cfg.Email.Sender)
	defer func() {
		if err := mailRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	mailUsecase := usecaseMail.NewMailUsecase(mailRepository)

	log.Printf("Connecting to Kafka at: %v", "kafka:9092")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := mailAppKafka.NewConsumer(
		[]string{"kafka:9092"},
		"on-email-sent-group",
		"balance_notification_topic",
		mailUsecase,
	)
	go func() {
		log.Println("Starting Kafka consumer...")
		consumer.Run(ctx)
	}()

	mux := mailAppHttp.SetupRoutes(mailUsecase)

	log.Fatal(http.ListenAndServe(":8036", mailMiddleware.CORS(mux)))
}
