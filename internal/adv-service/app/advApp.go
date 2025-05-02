package advApp

import (
	"log"
	"net/http"
	configs "retarget/configs"
	advAppHttp "retarget/internal/adv-service/controller/http"
	advMiddleware "retarget/internal/adv-service/controller/http/middleware"
	repoAdv "retarget/internal/adv-service/repo/adv"
	repoSlot "retarget/internal/adv-service/repo/slot"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"
	authenticate "retarget/pkg/middleware/auth"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	advRepository := repoAdv.NewAdvRepository("localhost", 9042, "slot_space", "cassandra", "12345678")

	slotRepository := repoSlot.NewSlotRepository("localhost", 9042, "slot_space", "cassandra", "12345678")
	defer slotRepository.Close()

	advUsecase := usecaseAdv.NewAdvUsecase(advRepository)

	slotUsecase := usecaseSlot.NewSlotUsecase(slotRepository)

	mux := advAppHttp.SetupRoutes(authenticator, advUsecase, slotUsecase)

	log.Fatal(http.ListenAndServe(":8032", advMiddleware.CORS(mux)))
}
