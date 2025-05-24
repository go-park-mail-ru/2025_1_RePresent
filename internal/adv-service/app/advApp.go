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
	pb "retarget/pkg/proto/banner"
	protoPayment "retarget/pkg/proto/payment"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	dsn := "clickhouse://user:123456@ReTargetClickHouse:9000/adv?dial_timeout=10s"
	advRepository := repoAdv.NewAdvRepository(cfg.Scylla.Host, cfg.Scylla.Port, cfg.Scylla.SlotKeyspace, cfg.Scylla.Username, cfg.Scylla.Password, dsn)

	slotRepository := repoSlot.NewSlotRepository(cfg.Scylla.Host, cfg.Scylla.Port, cfg.Scylla.SlotKeyspace, cfg.Scylla.Username, cfg.Scylla.Password)
	defer slotRepository.Close()

	conn, err := grpc.NewClient("ReTargetApiBanner:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	connPayment, err := grpc.NewClient("ReTargetApiPayment:8054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connPayment.Close()

	cBanner := pb.NewBannerServiceClient(conn)
	cPayment := protoPayment.NewPaymentServiceClient(connPayment)

	advUsecase := usecaseAdv.NewAdvUsecase(advRepository, cBanner, cPayment, slotRepository)

	slotUsecase := usecaseSlot.NewSlotUsecase(slotRepository)

	mux := advAppHttp.SetupRoutes(authenticator, advUsecase, slotUsecase)

	log.Fatal(http.ListenAndServe(":8032", advMiddleware.CORS(mux)))
}
