package controller

import (
	handlerCsat "retarget/internal/csat-serivce/controller/http/csat"
	usecaseCsat "retarget/internal/csat-serivce/usecase"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, csatUsecase *usecaseCsat.CsatUsecase) *mux.Router {
	r := mux.NewRouter()

	csatRoutes := handlerCsat.SetupCsatRoutes(authenticator, csatUsecase)
	r.PathPrefix("/api/v1/csat/").Handler(csatRoutes)

	return r
}
