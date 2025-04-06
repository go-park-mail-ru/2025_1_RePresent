package mail

import (
	"net/http"
	usecaseMail "retarget/internal/mail-service/usecase/mail"

	"github.com/gorilla/mux"
)

type MailController struct {
	mailUsecase *usecaseMail.MailUsecase
}

func NewMailController(mailUsecase *usecaseMail.MailUsecase) *MailController {
	return &MailController{mailUsecase: mailUsecase}
}

func SetupMailRoutes(mailUsecase *usecaseMail.MailUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	mailController := NewMailController(mailUsecase)

	muxRouter.Handle("/api/v1/mail/send-register-code", http.HandlerFunc(mailController.SendRegisterCodeHandler))
	// muxRouter.Handle("/api/v1/mail/send-recovery-code", http.HandlerFunc(mailController.SendRecoveryCodeHandler))
	// muxRouter.Handle("/api/v1/mail/send-password-reset-code", http.HandlerFunc(mailController.SendEditPasswordCodeHandler))

	return muxRouter
}
