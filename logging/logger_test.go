package logging

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Costiss/fs-backup/config"
)

func TestGetLogFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testlogging")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with a log file specified
	logFilePath := filepath.Join(tmpDir, "test.log")
	cfgWithLogFile := &config.Config{
		Backup: config.BackupConfig{
			LogFile: logFilePath,
		},
	}

	logFile := getLogFile(cfgWithLogFile)
	if logFile == nil {
		t.Errorf("Expected a log file, but got nil")
	}
	logFile.Close()

	// Verify that the log file was created
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFilePath)
	}

	// Test with no log file specified
	cfgWithoutLogFile := &config.Config{}
	logFile = getLogFile(cfgWithoutLogFile)
	if logFile != nil {
		t.Errorf("Expected nil, but got a log file")
	}
}
