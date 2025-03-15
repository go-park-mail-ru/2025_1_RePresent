package http

import (
	md "retarget/internal/controller/http/middleware"
	handlerMail "retarget/internal/mail-service/controller/http/mail"
	usecaseMail "retarget/internal/mail-service/usecase/mail"

	"github.com/gorilla/mux"
)

func SetupRoutes(mailUsecase *usecaseMail.MailUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(md.ErrorMiddleware)

	mailRoutes := handlerMail.SetupMailRoutes(mailUsecase)
	r.PathPrefix("/api/v1/mail/").Handler(mailRoutes)

	return r
}
