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
	pb "retarget/pkg/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	advRepository := repoAdv.NewAdvRepository("ReTargetScylla", 9042, "slot_space", "cassandra", "12345678")

	slotRepository := repoSlot.NewSlotRepository("ReTargetScylla", 9042, "slot_space", "cassandra", "12345678")
	defer slotRepository.Close()

	conn, err := grpc.Dial("ReTargetApiBanner:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBannerServiceClient(conn)

	advUsecase := usecaseAdv.NewAdvUsecase(advRepository, c)

	slotUsecase := usecaseSlot.NewSlotUsecase(slotRepository)

	mux := advAppHttp.SetupRoutes(authenticator, advUsecase, slotUsecase)

	log.Fatal(http.ListenAndServe(":8032", advMiddleware.CORS(mux)))
}
