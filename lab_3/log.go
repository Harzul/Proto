package main

import (
	"log"
	"os"
)

func initLogger() (*os.File, *log.Logger) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger := log.New(file, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)

	return file, logger
}
