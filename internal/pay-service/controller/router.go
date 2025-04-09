package http

import (
	handlerPayment "retarget/internal/pay-service/controller/http"
	usecasePayment "retarget/internal/pay-service/usecase"
	logs "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(paymentUsecase *usecasePayment.PaymentUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.LogMiddleware)

	paymentRoutes := handlerPayment.SetupPaymentRoutes(paymentUsecase)
	r.PathPrefix("/api/v1/payment/").Handler(paymentRoutes)

	return r
}
