package backup

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	FilePath  string
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

func UploadToS3(cfg S3Config) error {
	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
		config.WithBaseEndpoint(cfg.Endpoint),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(client)

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.FilePath),
		Body:   file,
	})

	return err
}
