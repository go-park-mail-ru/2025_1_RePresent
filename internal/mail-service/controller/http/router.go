package http

import (
	handlerMail "retarget/internal/mail-service/controller/http/mail"
	usecaseMail "retarget/internal/mail-service/usecase/mail"
	middleware "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(mailUsecase *usecaseMail.MailUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LogMiddleware)

	mailRoutes := handlerMail.SetupMailRoutes(mailUsecase)
	r.PathPrefix("/api/v1/mail/").Handler(mailRoutes)

	return r
}
