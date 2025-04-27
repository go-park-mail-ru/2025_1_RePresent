package main

import (
	"log"
	"os"
	"os/signal"
	"retarget/configs"
	csatApp "retarget/internal/csat-service/app"
	"syscall"
)

func main() {
	cfg, err := configs.LoadConfigs("configs/auth-redis.yml")
	if err != nil {
		log.Fatal(err)
	}
	go csatApp.Run(cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
