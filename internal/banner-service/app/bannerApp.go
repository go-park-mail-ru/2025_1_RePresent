package authApp

import (
	"fmt"
	"log"
	"net/http"
	configs "retarget/configs"
	controller "retarget/internal/banner-service/controller"
	middleware "retarget/internal/banner-service/controller/http/middleware"
	"retarget/internal/banner-service/repo"
	"retarget/internal/banner-service/service"
	authenticate "retarget/pkg/middleware/auth"

	server "retarget/internal/banner-service/grpc"
	usecase "retarget/internal/banner-service/usecase"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	gigaChatService := service.NewGigaChatService(logger, cfg.GigaChat.AuthKey, cfg.GigaChat.ClientID)
	imageRepository := repo.NewBannerImageRepository(cfg.Minio.EndPoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccesKey, cfg.Minio.Token, false, "image")
	bannerRepository := repo.NewBannerRepository(cfg.Database.ConnectionString("d"), logger, gigaChatService)
	defer func() {
		if err := bannerRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	image := usecase.NewBannerImageUsecase(imageRepository)
	banner := usecase.NewBannerUsecase(bannerRepository)

	mux := controller.SetupRoutes(authenticator, banner, image)

	errChan := make(chan error)

	// Запуск gRPC-сервера в горутине
	go func() {
		log.Println("Starting gRPC server...")
		server.RunGRPCServer(*banner)
	}()

	// Запуск HTTP-сервера в горутине
	go func() {
		log.Println("Starting HTTP server on :8024...")
		if err := http.ListenAndServe(":8024", middleware.CORS(mux)); err != nil {
			errChan <- fmt.Errorf("HTTP server failed: %v", err)
		}
	}()

	// Ожидание ошибки из любого сервера
	err = <-errChan
	log.Fatalf("Server error: %v", err)

}
