package mailApp

import (
	"log"
	"net/http"
	"retarget/configs"
	mailAppHttp "retarget/internal/mail-service/controller/http"
	mailMiddleware "retarget/internal/mail-service/controller/http/middleware"
	repoMail "retarget/internal/mail-service/repo/mail"
	usecaseMail "retarget/internal/mail-service/usecase/mail"
)

func Run(cfg *configs.Config) {
	mailRepository := repoMail.NewMailRepository(cfg.Email.SmtpServer, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, cfg.Email.Sender)
	defer func() {
		if err := mailRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	mailUsecase := usecaseMail.NewMailUsecase(mailRepository)

	mux := mailAppHttp.SetupRoutes(mailUsecase)

	log.Fatal(http.ListenAndServe(":8025", mailMiddleware.CORS(mux)))
}
