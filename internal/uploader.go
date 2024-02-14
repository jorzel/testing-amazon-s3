package internal

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

// Uploader defines the interface for uploading files.
type Uploader interface {
	UploadFile(bucketName, objectName, filePath string) error
}

// MinIOUploader implements the Uploader interface for MinIO.
type MinIOUploader struct {
	Client *minio.Client
}

// UploadFile uploads a file to the specified MinIO bucket.
func (u *MinIOUploader) UploadFile(ctx context.Context, bucketName string, filePath string) error {
	objectName := filepath.Base(filePath)
	_, err := u.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
