package advApp

import (
	"log"
	"net/http"
	"retarget/configs"
	advAppHttp "retarget/internal/adv-service/controller/http"
	advMiddleware "retarget/internal/adv-service/controller/http/middleware"
	repoAdv "retarget/internal/adv-service/repo/adv"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
)

func Run(cfg *configs.Config) {
	advRepository := repoAdv.NewAdvRepository("localhost", 9042, "link_space", "cassandra", "12345678")

	advUsecase := usecaseAdv.NewAdvUsecase(advRepository)

	mux := advAppHttp.SetupRoutes(advUsecase)

	log.Fatal(http.ListenAndServe(":8032", advMiddleware.CORS(mux)))
}
