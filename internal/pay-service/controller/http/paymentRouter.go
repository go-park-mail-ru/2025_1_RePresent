package payment

import (
	"net/http"
	payment "retarget/internal/pay-service/usecase"

	"github.com/gorilla/mux"
)

type PaymentController struct {
	PaymentUsecase *payment.PaymentUsecase
}

func NewPaymentController(PaymentUsecase *payment.PaymentUsecase) *PaymentController {
	return &PaymentController{PaymentUsecase: PaymentUsecase}
}

func SetupPaymentRoutes(PaymentUsecase *payment.PaymentUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	PaymentController := NewPaymentController(PaymentUsecase)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/accounts/{accountid}/balance", http.HandlerFunc(PaymentController.GetUserBalance))
	muxRouter.Handle("/api/v1/payment/accounts/{accountid}/topup", http.HandlerFunc(PaymentController.TopUpAccount))
	muxRouter.Handle("/api/v1/payment/transactions/{transactionid}", http.HandlerFunc(PaymentController.GetTransactionByID))
	//muxRouter.Handle("/api/v1/payment/transactions/{transactionid}/confirm", http.HandlerFunc(skibidi))

	return muxRouter
}
