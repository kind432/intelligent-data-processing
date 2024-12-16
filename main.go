package main

import (
	"intelligent-data-processing/internal/config"
	"intelligent-data-processing/internal/device"
	"intelligent-data-processing/internal/handler"
	"intelligent-data-processing/internal/mqtt"
	"intelligent-data-processing/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	loggers := logger.InitLogger()
	h := handler.NewHandler(loggers)
	client := mqtt.InitMQTTClient(h, loggers)
	defer func() {
		device.DisconnectAllDevices(h, client)
		client.Disconnect(250)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
