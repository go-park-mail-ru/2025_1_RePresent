package csat

import (
	"github.com/gorilla/mux"
	"net/http"
	csat "retarget/internal/csat-serivce/usecase"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"
)

type CsatController struct {
	csatUsecase *csat.CsatUsecase
}

func NewCsatController(csatUsecase *csat.CsatUsecase) CsatController {
	return CsatController{csatUsecase: csatUsecase}
}

func SetupCsatRoutes(authenticator *authenticate.Authenticator, csatUsecase *csat.CsatUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	csatController := NewCsatController(csatUsecase)

	muxRouter.Handle("/api/v1/csat/show/{page_id:[0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.ShowQuestionByPageID)))).Methods("GET")
	muxRouter.Handle("/api/v1/csat/send", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.SendReview)))).Methods("POST")

	return muxRouter
}
