package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	S3     S3Config     `yaml:"s3"`
	Backup BackupConfig `yaml:"backup"`
}

type S3Config struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type BackupConfig struct {
	Directories []string `yaml:"directories"`
	LogFile     string   `yaml:"log_file,omitempty"`
	Schedule    string   `yaml:"schedule"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
