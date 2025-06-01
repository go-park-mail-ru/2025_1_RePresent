package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"retarget/configs"
	mailApp "retarget/internal/mail-service/app"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("не удалось инициализировать логгер: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()
	sugar := logger.Sugar()
	cfg, err := configs.LoadConfigs()
	if err != nil {
		log.Fatal(err)
	}
	go mailApp.Run(cfg, sugar)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
