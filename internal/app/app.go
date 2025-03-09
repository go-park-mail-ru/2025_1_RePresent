package app

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"retarget/configs"
	appHttp "retarget/internal/controller/http"
	"retarget/internal/controller/http/middleware"
	"retarget/internal/repo"
	"retarget/internal/usecase/auth"
	"retarget/internal/usecase/banner"
)

func Run(cfg *configs.Config) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode))
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repo.NewUserRepository(db)
	authUsecase := auth.NewAuthUsecase(userRepository)

	bannerRepository := repo.NewBannerRepository(db)
	bannerUsecase := banner.NewBannerUsecase(bannerRepository)

	mux := appHttp.SetupRoutes(authUsecase, bannerUsecase)
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(mux)))
}
