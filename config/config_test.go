package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testconfig")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	configFile, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Create a dummy config struct
	config := Config{
		S3: S3Config{
			Bucket:    "test-bucket",
			Region:    "us-east-1",
			Endpoint:  "http://localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
		},
		Backup: BackupConfig{
			Directories: []string{"/tmp/test1", "/tmp/test2"},
			Schedule:    "@daily",
			Database: []PGConfig{
				{
					Host:     "localhost",
					Port:     5432,
					User:     "testuser",
					Password: "testpassword",
				},
			},
		},
	}

	// Marshal the config struct to YAML and write it to the file
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		t.Fatalf("Failed to marshal config to YAML: %v", err)
	}
	_, err = configFile.Write(yamlData)
	if err != nil {
		t.Fatalf("Failed to write YAML to config file: %v", err)
	}
	configFile.Close()

	// Set the environment variable for the password
	passwordEnvVar := "PG_PASSWORD_testuser"
	os.Setenv(passwordEnvVar, "newpassword")
	defer os.Unsetenv(passwordEnvVar)

	// Load the config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the config
	if loadedConfig.S3.Bucket != config.S3.Bucket {
		t.Errorf("Expected S3 bucket %s, got %s", config.S3.Bucket, loadedConfig.S3.Bucket)
	}
	if loadedConfig.Backup.Directories[0] != config.Backup.Directories[0] {
		t.Errorf("Expected directory %s, got %s", config.Backup.Directories[0], loadedConfig.Backup.Directories[0])
	}
	if loadedConfig.Backup.Database[0].Password != "newpassword" {
		t.Errorf("Expected password to be overridden to 'newpassword', got '%s'", loadedConfig.Backup.Database[0].Password)
	}
}
