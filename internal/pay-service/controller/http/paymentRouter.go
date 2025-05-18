package payment

import (
	"net/http"
	payment "retarget/internal/pay-service/usecase"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type PaymentController struct {
	PaymentUsecase *payment.PaymentUsecase
}

func NewPaymentController(PaymentUsecase *payment.PaymentUsecase) *PaymentController {
	return &PaymentController{PaymentUsecase: PaymentUsecase}
}

func SetupPaymentRoutes(authenticator *authenticate.Authenticator, PaymentUsecase *payment.PaymentUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	PaymentController := NewPaymentController(PaymentUsecase)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/payment/balance", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(PaymentController.GetUserBalance))))
	muxRouter.Handle("/api/v1/payment/accounts/topup", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(PaymentController.TopUpAccount))))
	muxRouter.Handle("/api/v1/payment/transactions/clicks", logger.LogMiddleware(http.HandlerFunc(PaymentController.RegUserActivity)))
	//muxRouter.Handle("/api/v1/payment/transactions/{transactionid}/confirm", http.HandlerFunc(skibidi))

	muxRouter.Handle("/api/v1/payment/transactions/{transactionid}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(PaymentController.GetTransactionByID))))
	muxRouter.Handle("/api/v1/payment/transactions", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(PaymentController.CreateTransaction)))).Methods("POST")

	return muxRouter
}
