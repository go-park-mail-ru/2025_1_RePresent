package controller

import (
	"log"
	handlerCsat "retarget/internal/csat-service/controller/http/csat"
	usecaseCsat "retarget/internal/csat-service/usecase/csat"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, csatUsecase *usecaseCsat.CsatUsecase) *mux.Router {
	r := mux.NewRouter()
	log.Printf("Router")

	csatRoutes := handlerCsat.SetupCsatRoutes(authenticator, csatUsecase)
	r.PathPrefix("/api/v1/csat/").Handler(csatRoutes)

	return r
}
