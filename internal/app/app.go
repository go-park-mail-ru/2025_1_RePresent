package app

import (
	"database/sql"

	"net/http"

	appHttp "RE/internal/controller/http"
	"RE/internal/controller/http/middleware"

	"RE/internal/repo"
	"RE/internal/usecase/auth"
)

func Run() {
	db, err := sql.Open("postgres", "user=postgres password=123456 dbname=testdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	userRepository := repo.NewUserRepository(db)
	authUsecase := auth.NewAuthUsecase(userRepository)
	mux := appHttp.SetupRoutes(authUsecase)
	http.ListenAndServe(":8080", middleware.CORS(mux))
}
