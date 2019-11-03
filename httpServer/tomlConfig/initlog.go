package tomlConfig

import (
	"fmt"
	"io"
	"log"
	"os"
)

func InitLogger() *log.Logger {
	// get log file
	logFile := conf.HttpServer.LogFile

	// open logFile if it exsit or create it if it not exsit
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("create of open logFile", logFile, "is bad", err)
		return nil
	}

	// init log.logger
	return log.New(io.MultiWriter(os.Stdout, file), "[NOS] ", log.Ldate|log.Ltime|log.Lshortfile)
}
