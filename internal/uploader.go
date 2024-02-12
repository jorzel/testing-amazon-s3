package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/minio/minio-go"
)

// Uploader defines the interface for uploading files.
type Uploader interface {
	UploadFile(bucketName, objectName, filePath string) error
}

// S3Uploader implements the Uploader interface for Amazon S3.
type S3Uploader struct {
	Service *s3.S3
}

// UploadFile uploads a file to the specified S3 bucket.
func (u *S3Uploader) UploadFile(bucketName, objectName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = u.Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	log.Println("File uploaded")
	return nil
}

// MinIOUploader implements the Uploader interface for MinIO.
type MinIOUploader struct {
	Client *minio.Client
}

// UploadFile uploads a file to the specified MinIO bucket.
func (u *MinIOUploader) UploadFile(ctx context.Context, bucketName string, objectName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = u.Client.PutObjectWithContext(ctx, bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
