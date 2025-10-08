package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var s3Client *s3.Client
var resource *dockertest.Resource

const (
	minioUser     = "minioadmin"
	minioPassword = "minioadmin"
	bucketName    = "test-bucket"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Cmd:        []string{"server", "/data"},
		Env: []string{
			"MINIO_ACCESS_KEY=" + minioUser,
			"MINIO_SECRET_KEY=" + minioPassword,
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	endpoint := fmt.Sprintf("http://localhost:%s", resource.GetPort("9000/tcp"))

	if err := pool.Retry(func() error {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		awsCfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioUser, minioPassword, "")),
		)
		if err != nil {
			return err
		}

		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})

		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestUploadToS3(t *testing.T) {
	// Create a temporary file to upload
	tmpDir, err := os.MkdirTemp("", "test-s3")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test-file.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	f.WriteString("hello world")
	f.Close()

	endpoint := fmt.Sprintf("http://localhost:%s", resource.GetPort("9000/tcp"))
	// Upload the file to the mock S3 server
	s3Cfg := S3Config{
		FilePath:  filePath,
		Bucket:    bucketName,
		Region:    "us-east-1",
		Endpoint:  endpoint,
		AccessKey: minioUser,
		SecretKey: minioPassword,
	}

	err = UploadToS3(s3Cfg)
	if err != nil {
		t.Fatalf("UploadToS3 failed: %v", err)
	}

	// Verify that the file was uploaded successfully
	_, err = s3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath.Base(filePath)),
	})
	if err != nil {
		t.Errorf("failed to find uploaded file in mock S3 server: %v", err)
	}
}
