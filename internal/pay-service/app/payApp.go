package mailApp

import (
	"log"
	"net/http"
	"retarget/configs"
	payAppHttp "retarget/internal/pay-service/controller"
	payMiddleware "retarget/internal/pay-service/controller/http/middleware"
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

	dsn := "clickhouse://user:123456@ReTargetClickHouse:9000/csat"
	payRepository := repoPay.NewPaymentRepository(cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode, dsn, logger)
	defer func() {
		if err := payRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	payUsecase := usecasePay.NewPayUsecase(payRepository)

	mux := payAppHttp.SetupRoutes(authenticator, payUsecase)

	log.Fatal(http.ListenAndServe(":8022", payMiddleware.CORS(mux)))
}
