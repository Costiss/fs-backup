package main

import (
	"flag"
	"fmt"

	"github.com/Costiss/fs-backup/backup"
	"github.com/Costiss/fs-backup/config"
	"github.com/Costiss/fs-backup/logging"
	"github.com/robfig/cron/v3"
)

var Version = "dev"

func main() {
	var (
		runOnce  = flag.Bool("run-once", false, "Run backup immediately and exit")
		showHelp = flag.Bool("help", false, "Show help message")
		version  = flag.Bool("version", false, "Show version and exit")
	)
	flag.Parse()

	if *showHelp {
		fmt.Println("Usage: fs-backup [options] [config_path]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Printf("fs-backup version: %s\n", Version)
		return
	}

	configPath := "/etc/fs-backup/config.yaml"
	if len(flag.Args()) > 0 {
		configPath = flag.Args()[0]
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

	if *runOnce {
		logger.Println("==============================================")
		logger.Println("         🚀 FS-BACKUP: Running Once 🚀        ")
		logger.Println("==============================================")
		backup.Run(cfg)
		return
	}

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
