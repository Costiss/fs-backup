package backup

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Costiss/fs-backup/config"
	"github.com/Costiss/fs-backup/logging"
)

func Run(cfg *config.Config) {
	logger := logging.GetLogger("backup_runner", cfg)

	logger.Println("==============================================")
	logger.Println("           	   Starting backup         	      ")
	logger.Println("==============================================")

	for idx, dir := range cfg.Backup.Directories {
		logger.Println("----------------------------------------------")
		logger.Printf(" [%d/%d] Backing up directory: %s\n", idx+1, len(cfg.Backup.Directories), dir)
		logger.Println("----------------------------------------------")

		tarpath := dir + ".tar.gz"
		err := compressDirectory(dir, tarpath)
		if err != nil {
			logger.Printf("  ✗ Error compressing directory: %v\n", err)
			logger.Println("----------------------------------------------")
			continue
		}

		logger.Printf("  ✓ Created archive: %s\n", tarpath)

		s3Cfg := S3Config{
			FilePath:  tarpath,
			Bucket:    cfg.S3.Bucket,
			Region:    cfg.S3.Region,
			Endpoint:  cfg.S3.Endpoint,
			AccessKey: cfg.S3.AccessKey,
			SecretKey: cfg.S3.SecretKey,
		}
		if err := UploadToS3(s3Cfg); err != nil {
			logger.Printf("  ✗ Error uploading to S3: %v\n", err)
			logger.Println("----------------------------------------------")
			os.Remove(tarpath)
			continue
		}

		logger.Printf("  ✓ Successfully backed up directory: %s\n", dir)
		logger.Println("----------------------------------------------")
		os.Remove(tarpath)
	}

	logger.Println("==============================================")
	logger.Println("           	  Backup completed                ")
	logger.Println("==============================================")
}

func compressDirectory(srcDir, destFile string) error {
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(srcDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(filepath.Dir(srcDir), file)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(tw, f)
			if err != nil {
				return err
			}
		}
		return nil
	})

}

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
