package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}

func New() Logger {
	logger := Logger{
		Info: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		Err:  log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime),
	}
	logger.Info.Print("Logger initialized")
	return logger
}
