package http

import (
	// handlerAdv "retarget/internal/adv-service/controller/http/adv"
	handlerAdv "retarget/internal/adv-service/controller/http/adv"
	handlerSlot "retarget/internal/adv-service/controller/http/slot"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"

	// middleware "retarget/pkg/middleware"
	"retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, advUsecase *usecaseAdv.AdvUsecase, slotUsecase *usecaseSlot.SlotUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LogMiddleware)

	advRoutes := handlerAdv.SetupAdvRoutes(authenticator, advUsecase, slotUsecase)
	r.PathPrefix("/api/v1/adv/").Handler(advRoutes)

	slotRoutes := handlerSlot.SetupSlotRoutes(authenticator, slotUsecase)
	r.PathPrefix("/api/v1/slot/").Handler(slotRoutes)

	return r
}
