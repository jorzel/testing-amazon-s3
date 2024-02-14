package main

import (
	"context"
	"log"
	"os"

	"amazon-s3-uploader/internal"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

type UploadSpec struct {
	accessKey  string
	secretKey  string
	bucketName string
	fileName   string
}

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "uploader",
	}

	cmd.AddCommand(
		uploadCmd(),
	)
	return cmd
}

func uploadCmd() *cobra.Command {
	var fileName string
	var bucketName string
	var accessKey string
	var secretKey string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload file to s3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName, err := cmd.Flags().GetString("filename")
			if err != nil {
				return err
			}
			bucketName, err := cmd.Flags().GetString("bucket")
			if err != nil {
				return err
			}
			accessKey, err := cmd.Flags().GetString("accesskey")
			if err != nil {
				return err
			}
			secretKey, err := cmd.Flags().GetString("secretkey")
			if err != nil {
				return err
			}
			return upload(UploadSpec{
				fileName:   fileName,
				bucketName: bucketName,
				accessKey:  accessKey,
				secretKey:  secretKey,
			})
		},
	}
	cmd.PersistentFlags().StringVarP(&fileName, "filename", "f", "", "Full file path to resource")
	cmd.MarkPersistentFlagRequired("filename")
	cmd.PersistentFlags().StringVarP(&accessKey, "accesskey", "a", "", "Amazon S3 Access Key")
	cmd.MarkPersistentFlagRequired("accesskey")
	cmd.PersistentFlags().StringVarP(&secretKey, "secretkey", "x", "", "Amazon S3 Secret Key")
	cmd.MarkPersistentFlagRequired("secretkey")
	cmd.PersistentFlags().StringVarP(&bucketName, "bucket", "b", "", "Amazon S3 Bucket Name")
	cmd.MarkPersistentFlagRequired("bucket")
	return cmd
}

func upload(spec UploadSpec) error {
	endpoint := "s3.amazonaws.com"

	// Initialize AWS session for S3Uploader
	s3Client, err := minio.New(
		endpoint,
		&minio.Options{Secure: true, Creds: credentials.NewStaticV4(spec.accessKey, spec.secretKey, "")},
	)
	if err != nil {
		return err
	}
	s3Uploader := &internal.MinIOUploader{Client: s3Client}

	err = s3Uploader.UploadFile(context.Background(), spec.bucketName, spec.fileName)
	if err != nil {
		return err
	}

	log.Printf("Upload of: %s to bucket: %s success", spec.fileName, spec.bucketName)
	return nil
}
