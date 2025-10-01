package logging

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Costiss/fs-backup/config"
)

func getLogFile(cfg *config.Config) *os.File {
	logFilePath := cfg.Backup.LogFile
	if logFilePath == "" {
		logFilePath = "/var/log/fs-backup.log"
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return logFile
}

func GetLogger(name string, cfg *config.Config) *log.Logger {
	logFile := getLogFile(cfg)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, fmt.Sprintf("%s: ", name), log.LstdFlags)
	return logger
}
