package adv

import (
	"net/http"
	AdvMiddleware "retarget/internal/adv-service/controller/http/middleware"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"

	"github.com/gorilla/mux"
)

type AdvController struct {
	advUsecase *usecaseAdv.AdvUsecase
}

func NewAdvController(advUsecase *usecaseAdv.AdvUsecase) *AdvController {
	return &AdvController{advUsecase: advUsecase}
}

func SetupAdvRoutes(advUsecase *usecaseAdv.AdvUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	advController := NewAdvController(advUsecase)

	advMiddleware := AdvMiddleware.AdvMiddleware(advUsecase)

	muxRouter.Handle("/api/v1/adv/link/generate", http.HandlerFunc(advController.GenerateLinkHandler))
	muxRouter.Handle("/api/v1/adv/link/my", http.HandlerFunc(advController.GetLinkHandler))

	muxRouter.Handle("/api/v1/adv/iframe/{link}", advMiddleware(http.HandlerFunc(advController.IframeHandler)))
	muxRouter.Handle("/api/v1/adv/metrics/{link}", advMiddleware(http.HandlerFunc(advController.MetricsHandler)))

	return muxRouter
}
