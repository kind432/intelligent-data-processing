package logger

import (
	"log"
	"os"
)

type Loggers struct {
	Info *log.Logger
	Err  *log.Logger
}

func InitLogger() (loggers Loggers) {
	//infoF, err := os.OpenFile(viper.GetString("logger.info"), os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	log.Fatalf("%s", err.Error())
	//}
	//defer func(infoF *os.File) {
	//	err := infoF.Close()
	//	if err != nil {
	//		loggers.Err.Fatalf("%s", err.Error())
	//	}
	//}(infoF)
	//errF, err := os.OpenFile(viper.GetString("logger.error"), os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	log.Fatalf("%s", err.Error())
	//}
	//defer func(errF *os.File) {
	//	err := errF.Close()
	//	if err != nil {
	//		loggers.Err.Fatalf("%s", err.Error())
	//	}
	//}(errF)
	//
	//loggers.Info = log.New(infoF, "[INFO]\t", log.Ldate|log.Ltime)
	//loggers.Err = log.New(errF, "[ERROR]\t", log.Ldate|log.Ltime)
	loggers.Info = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	loggers.Err = log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime)
	loggers.Info.Print("Logger initialized.")
	return loggers
}
