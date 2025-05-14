package mailApp

import (
	"log"
	"net/http"
	"retarget/configs"
	payAppHttp "retarget/internal/pay-service/controller"
	payMiddleware "retarget/internal/pay-service/controller/http/middleware"
	server "retarget/internal/pay-service/grpc"
	repoPay "retarget/internal/pay-service/repo"
	repoAttempt "retarget/internal/pay-service/repo/attempt"
	repoNotice "retarget/internal/pay-service/repo/notice"
	usecasePay "retarget/internal/pay-service/usecase"
	authenticate "retarget/pkg/middleware/auth"
	"time"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	payRepository := repoPay.NewPaymentRepository(cfg.Database.UsernameAuth, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode, logger)
	defer func() {
		if err := payRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()
	noticeRepository := repoNotice.NewNoticeRepository([]string{"kafka:9092"}, "notifications", logger) // balance_notification_topic
	if noticeRepository == nil {
		logger.Fatal("failed to initialize NoticeRepository")
	}
	defer noticeRepository.Close()
	attemptRepository := repoAttempt.NewAttemptRepository(cfg.AttemptRedis.EndPoint, cfg.AttemptRedis.Password, cfg.AttemptRedis.Database, 24*time.Hour, cfg.AttemptRedis.Attempts)

	payUsecase := usecasePay.NewPayUsecase(logger, payRepository, noticeRepository, attemptRepository)

	go func() {
		log.Println("Starting gRPC server...")
		server.RunGRPCServer(*payUsecase)
	}()

	mux := payAppHttp.SetupRoutes(authenticator, payUsecase)

	log.Fatal(http.ListenAndServe(":8022", payMiddleware.CORS(mux)))
}
