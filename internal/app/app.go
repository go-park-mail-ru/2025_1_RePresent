package app

import (
	"database/sql"
	"log"

	"net/http"

	appHttp "RE/internal/controller/http"
	"RE/internal/controller/http/middleware"

	"RE/internal/repo"
	"RE/internal/usecase/auth"
)

func Run() {
	db, err := sql.Open("postgres", "user=postgres password=123456 dbname=test host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	userRepository := repo.NewUserRepository(db)
	authUsecase := auth.NewAuthUsecase(userRepository)
	mux := appHttp.SetupRoutes(authUsecase)
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(mux)))
}
