package adv

import (
	"net/http"
	AdvMiddleware "retarget/internal/adv-service/controller/http/middleware"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"
	authenticate "retarget/pkg/middleware/auth"
	pb "retarget/pkg/proto/banner"
	"time"

	"github.com/gorilla/mux"
)

// AdvUsecaseInterface определяет интерфейс для случая использования рекламы
type AdvUsecaseInterface interface {
	WriteMetric(bannerID int, slotLink string, action string) error
	GetIframe(secretLink string) (*pb.Banner, error)
	GetSlotMetric(slotID string, activity string, userID int, from, to time.Time) (interface{}, error)
	GetSlotCTR(slotID string, activity string, userID int, from, to time.Time) (interface{}, error)
	GetSlotRevenue(slotID string, activity string, userID int, from, to time.Time) (interface{}, error)
	GetSlotAVGPrice(slotID string, activity string, userID int, from, to time.Time) (interface{}, error)
	GetBannerMetric(bannerID int, activity string, userID int, from, to time.Time) (interface{}, error)
	GetBannerCTR(bannerID int, activity string, userID int, from, to time.Time) (interface{}, error)
	GetBannerExpenses(bannerID int, activity string, userID int, from, to time.Time) (interface{}, error)
}

type AdvController struct {
	advUsecase *usecaseAdv.AdvUsecase
}

func NewAdvController(advUsecase *usecaseAdv.AdvUsecase) *AdvController {
	return &AdvController{advUsecase: advUsecase}
}

func SetupAdvRoutes(authenticator *authenticate.Authenticator, advUsecase *usecaseAdv.AdvUsecase, slotUsecase *usecaseSlot.SlotUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	advController := NewAdvController(advUsecase)

	advMiddleware := AdvMiddleware.LinkMiddleware(slotUsecase)

	muxRouter.Handle("/api/v1/adv/iframe/{link}", advMiddleware(http.HandlerFunc(advController.IframeHandler))).Methods("GET")
	muxRouter.Handle("/api/v1/adv/metrics/", http.HandlerFunc(advController.MetricsHandler)).Methods("GET")
	muxRouter.Handle("/api/v1/adv/my-metrics", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(advController.MyMetricsHandler))).Methods("GET")

	return muxRouter
}
