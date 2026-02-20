package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"intelligent-data-processing/internal/config"
	"intelligent-data-processing/internal/repository"
	"intelligent-data-processing/internal/service"
	"intelligent-data-processing/internal/transport/mqtt"
	"intelligent-data-processing/pkg/logger"
)

func main() {
	if err := config.New(); err != nil {
		log.Fatalf("Config error: %s", err.Error())
	}

	loggers := logger.New()

	repo := repository.NewStorage()
	proc := service.NewProcessor(repo)
	handler := mqtt.NewHandler(loggers, repo, proc)

	client := mqtt.NewClient(loggers, handler)
	defer client.Disconnect(250)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	loggers.Info.Println("Shutting down gracefully...")
}
