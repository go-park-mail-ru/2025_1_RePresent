package adv

import (
	"net/http"
	AdvMiddleware "retarget/internal/adv-service/controller/http/middleware"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

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
	muxRouter.Handle("/api/v1/adv/metrics/{link}", http.HandlerFunc(advController.MetricsHandler))

	return muxRouter
}
