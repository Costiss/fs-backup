package main

import (
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/user/fs-backup/backup"
	"github.com/user/fs-backup/config"
	"github.com/user/fs-backup/logging"
)

func main() {
	configPath := "/etc/fs-backup/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error: Failed to load configuration: %v", err)
		return
	}

	if cfg.S3.AccessKey == "" || cfg.S3.SecretKey == "" {
		fmt.Println("Error: S3 access key or secret key is missing in the configuration.")
		return
	}

	logger := logging.GetLogger("main", cfg)

	c := cron.New()
	_, err = c.AddFunc(cfg.Backup.Schedule, func() {
		backup.Run(cfg)
	})
	if err != nil {
		logger.Printf("Error: Failed to schedule backup: %v", err)
		return
	}

	c.Start()
	logger.Println("==============================================")
	logger.Println("         🚀 FS-BACKUP: Started Running 🚀      ")
	logger.Println("==============================================")
	logger.Printf("  🕒 Cron schedule    : %s", cfg.Backup.Schedule)
	logger.Printf("  🪣 S3 bucket        : %s", cfg.S3.Bucket)
	logger.Printf("  🌍 S3 region        : %s", cfg.S3.Region)
	logger.Printf("  🔗 S3 endpoint      : %s", cfg.S3.Endpoint)
	logger.Printf("  📁 Backup dirs      : %v", cfg.Backup.Directories)
	logger.Printf("  📝 Log file         : %s", cfg.Backup.LogFile)
	logger.Println("==============================================")

	select {}
}
