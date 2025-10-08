package backup

import (
	"github.com/Costiss/fs-backup/config"
	"github.com/Costiss/fs-backup/logging"
)

func Run(cfg *config.Config) {
	logger := logging.GetLogger("backup_runner", cfg)

	logger.Println("==============================================")
	logger.Println("           	   Starting backup         	      ")
	logger.Println("==============================================")

	if len(cfg.Backup.Directories) != 0 {
		DoFsBackup(cfg)
	}

	logger.Println("==============================================")
	logger.Println("           	  Backup completed                ")
	logger.Println("==============================================")
}
