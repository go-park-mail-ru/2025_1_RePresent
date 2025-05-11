package main

import (
	"log"
	"os"
	"os/signal"
	"retarget/configs"
	app "retarget/internal/pay-service/app"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("не удалось инициализировать логгер: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := configs.LoadConfigs()
	if err != nil {
		sugar.Fatal(err)
	}
	go app.Run(cfg, sugar)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
