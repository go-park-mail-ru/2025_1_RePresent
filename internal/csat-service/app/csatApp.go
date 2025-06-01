package csatApp

import (
	"log"
	"net/http"
	"retarget/configs"
	csatAppHttp "retarget/internal/csat-service/controller/http"
	csatMiddleware "retarget/internal/csat-service/controller/http/middleware"
	repoCsat "retarget/internal/csat-service/repo/csat"
	usecaseCsat "retarget/internal/csat-service/usecase/csat"
	authenticate "retarget/pkg/middleware/auth"
)

func Run(cfg *configs.Config) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	dsn := "clickhouse://user:123456@ReTargetClickHouse:9000/csat?dial_timeout=10s"
	csatRepository := repoCsat.NewCsatRepository(dsn)
	defer func() {
		if err := csatRepository.CloseConnection(); err != nil {
			log.Printf("error closing CSAT repository: %v", err)
		}
	}()
	log.Printf("app")

	csatUsecase := usecaseCsat.NewCsatUsecase(csatRepository)

	mux := csatAppHttp.SetupRoutes(authenticator, csatUsecase)

	log.Fatal(http.ListenAndServe(":8035", csatMiddleware.CORS(mux)))
}
