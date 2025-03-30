package http

import (
	handlerAdv "retarget/internal/adv-service/controller/http/adv"
	usecaseAdv "retarget/internal/adv-service/usecase/adv"
	middleware "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(advUsecase *usecaseAdv.AdvUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LogMiddleware)

	advRoutes := handlerAdv.SetupAdvRoutes(advUsecase)
	r.PathPrefix("/api/v1/adv/").Handler(advRoutes)

	return r
}
