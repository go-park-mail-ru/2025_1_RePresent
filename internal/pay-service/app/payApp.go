package mailApp

import (
	"log"
	"net/http"
	"retarget/configs"
	payAppHttp "retarget/internal/pay-service/controller"
	payMiddleware "retarget/internal/pay-service/controller/http/middleware"
	server "retarget/internal/pay-service/grpc"
	repoPay "retarget/internal/pay-service/repo"
	usecasePay "retarget/internal/pay-service/usecase"
	authenticate "retarget/pkg/middleware/auth"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	payRepository := repoPay.NewPaymentRepository(cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode, logger)
	defer func() {
		if err := payRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	payUsecase := usecasePay.NewPayUsecase(payRepository)

	go func() {
		log.Println("Starting gRPC server...")
		server.RunGRPCServer(*payUsecase)
	}()

	mux := payAppHttp.SetupRoutes(authenticator, payUsecase)

	log.Fatal(http.ListenAndServe(":8022", payMiddleware.CORS(mux)))
}
