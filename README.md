## testing-amazon-s3
How to test amazon s3 upload using MinIO client.

### FileUploader
```go
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

```

### Integration tests with MinIO and Testcontainers
```go
package internal

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startMinioServer(ctx context.Context, accessKey string, secretKey string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"MINIO_ACCESS_KEY": accessKey,
			"MINIO_SECRET_KEY": secretKey,
		},
		Image: "minio/minio",

		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForHTTP("/minio/health/live").WithPort("9000"),
		Cmd:          []string{"server", "/data"},
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func terminateMinioServer(ctx context.Context, container testcontainers.Container) error {
	return container.Terminate(ctx)
}

func TestUploaderWithMinioServer(t *testing.T) {
	ctx := context.Background()

	accessKey := "testAccessKey"
	secretKey := "testSecretKey"
	bucketName := "testbucket"

	minioServer, err := startMinioServer(ctx, accessKey, secretKey)
	require.NoError(t, err)
	defer terminateMinioServer(ctx, minioServer)

	endpoint, err := minioServer.Endpoint(ctx, "")
	if err != nil {
		t.Fatal("cannot setup minio endpoint: %w", err)
	}
	minioClient, err := minio.New(
		endpoint,
		&minio.Options{
			Secure: false,
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		},
	)
	if err != nil {
		t.Fatal("cannot initialize minio client: %w", err)
	}
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		t.Fatal("cannot create minio bucket: %w", err)
	}
	fileName := "test.png"
	basePath := t.TempDir()
	filePath := filepath.Join(basePath, fileName)
	createFile(t, filePath)
	uploader := MinIOUploader{Client: minioClient}

	err = uploader.UploadFile(ctx, bucketName, filePath)
	require.NoError(t, err)

	_, err = minioClient.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	require.NoError(t, err)

}

func createFile(t *testing.T, filePath string) *os.File {
	t.Helper()
	fp, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("add file(%q): %s", filePath, err)
	}

	return fp
}
```
