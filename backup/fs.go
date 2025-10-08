package backup

import (
	"os"

	"github.com/Costiss/fs-backup/config"
	"github.com/Costiss/fs-backup/logging"
)

func DoFsBackup(cfg *config.Config) {
	logger := logging.GetLogger("fs_backup", cfg)

	for idx, dir := range cfg.Backup.Directories {
		logger.Println("----------------------------------------------")
		logger.Printf(" [%d/%d] Backing up directory: %s\n", idx+1, len(cfg.Backup.Directories), dir)
		logger.Println("----------------------------------------------")

		tarpath := dir + ".tar.gz"
		encryptedTarpath := tarpath + ".gpg"
		err := compressDirectory(dir, tarpath)
		if err != nil {
			logger.Printf("  ✗ Error compressing directory: %v\n", err)
			logger.Println("----------------------------------------------")
			continue
		}

		finalPath := tarpath
		if cfg.Backup.GpgEncryptPassword != "" {
			if err := encryptFileWithGPG(tarpath, encryptedTarpath, cfg.Backup.GpgEncryptPassword); err != nil {
				logger.Printf("  ✗ Error encrypting archive: %v\n", err)
				logger.Println("----------------------------------------------")
				continue
			}
			finalPath = encryptedTarpath
		}
		defer os.Remove(finalPath)

		logger.Printf("  ✓ Created archive: %s\n", finalPath)

		s3Cfg := S3Config{
			FilePath:  finalPath,
			Bucket:    cfg.S3.Bucket,
			Region:    cfg.S3.Region,
			Endpoint:  cfg.S3.Endpoint,
			AccessKey: cfg.S3.AccessKey,
			SecretKey: cfg.S3.SecretKey,
		}
		if err := UploadToS3(s3Cfg); err != nil {
			logger.Printf("  ✗ Error uploading to S3: %v\n", err)
			logger.Println("----------------------------------------------")
			continue
		}

		logger.Printf("  ✓ Successfully backed up directory: %s\n", dir)
		logger.Println("----------------------------------------------")
	}

}
