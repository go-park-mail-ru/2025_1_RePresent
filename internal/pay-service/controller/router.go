package http

import (
	handlerPayment "retarget/internal/pay-service/controller/http"
	usecasePayment "retarget/internal/pay-service/usecase"
	authenticate "retarget/pkg/middleware/auth"

	logs "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, paymentUsecase *usecasePayment.PaymentUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.LogMiddleware)

	paymentRoutes := handlerPayment.SetupPaymentRoutes(authenticator, paymentUsecase)
	r.PathPrefix("/api/v1/payment/").Handler(paymentRoutes)

	return r
}
