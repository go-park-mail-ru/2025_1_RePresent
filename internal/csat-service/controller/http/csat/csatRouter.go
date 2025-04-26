package csat

import (
	"net/http"
	csat "retarget/internal/csat-service/usecase/csat"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
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

	muxRouter.Handle("/api/v1/csat/show/{page_id:[a-zA-Z0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.ShowQuestionByPageID)))).Methods("GET")
	muxRouter.Handle("/api/v1/csat/my-reviews", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.ShowReviewsByUserID)))).Methods("GET")
	// muxRouter.Handle("/api/v1/csat/reviews", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.ShowAllReviews)))).Methods("GET")
	muxRouter.Handle("/api/v1/csat/show/iframe/{page_id:[a-zA-Z0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.ShowQuestionIFrameByPageID)))).Methods("GET")
	muxRouter.Handle("/api/v1/csat/send", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.SendReview)))).Methods("POST")

	// TODO Аналитика
	// muxRouter.Handle("/api/v1/csat/put-questions", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(csatController.SendReview)))).Methods("POST")

	return muxRouter
}
