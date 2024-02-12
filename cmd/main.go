package main

import (
	"context"
	"log"

	"amazon-s3-uploader/internal"

	"github.com/aws/aws-sdk-go/aws"
	awsCred "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/minio/minio-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	accessKey := "key"
	secretKey := "secret"
	region := "eu-central-1"
	bucketName := "jorzelbucket"

	// Initialize AWS session for S3Uploader
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: awsCred.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Fatalf("failed to initialize AWS session: %v", err)
	}
	svc := s3.New(sess)
	s3Uploader := &internal.S3Uploader{Service: svc}

	service := internal.UploadService{
		Uploader: s3Uploader,
	}
	service.Upload(context.Background(), "./")
	// Upload file using S3Uploader
	// err = s3Uploader.UploadFile(bucketName, "example.md", "./README.md")
	// if err != nil {
	// 	log.Fatalf("failed to upload file to S3: %v", err)
	// }
	//https://minioserver.example.net
	// Initialize MinIO client for MinIOUploader

	ctx := context.Background()
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
	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer minioContainer.Terminate(ctx)

	// minioContainer, err := testMinio.RunContainer(ctx, testcontainers.WithImage("minio/minio"))
	// if err != nil {
	// 	log.Fatalf("failed to start container: %s", err)
	// }

	// // Clean up the container
	// defer func() {
	// 	if err := minioContainer.Terminate(ctx); err != nil {
	// 		log.Fatalf("failed to terminate container: %s", err)
	// 	}
	// }()

	// Get MinIO container IP and port
	endpoint, err := minioContainer.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	// Initialize MinIO client
	minioClient, err := minio.New(
		endpoint,
		accessKey,
		secretKey,
		false, // Change to true if MinIO is configured with TLS
	)
	if err != nil {
		panic(err)
	}

	// Create a bucket for testing
	err = minioClient.MakeBucket(bucketName, region)
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}
	minioUploader := &internal.MinIOUploader{Client: minioClient}

	// Upload file using MinIOUploader
	err = minioUploader.UploadFile(ctx, bucketName, "example.md", "./README.md")
	if err != nil {
		log.Fatalf("failed to upload file to MinIO: %v", err)
	}
	log.Println("Upload success")
}
