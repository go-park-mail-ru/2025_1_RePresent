package app

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	appHttp "RE/internal/controller/http"
	"RE/internal/controller/http/middleware"
	"RE/internal/repo"
	"RE/internal/usecase/auth"

	"RE/configs"
)

func Run(cfg *configs.Config) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode))
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repo.NewUserRepository(db)
	authUsecase := auth.NewAuthUsecase(userRepository)
	mux := appHttp.SetupRoutes(authUsecase)
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(mux)))
}
