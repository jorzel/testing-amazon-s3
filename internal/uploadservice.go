package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type UploadService struct {
	Uploader Uploader
}

func (u UploadService) Upload(ctx context.Context, sourcePath string) {
	getFiles(sourcePath)

}

func getFiles(sourcePath string) []string {
	allowedExtension := ".md"
	foundFiles, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	var filePaths []string
	path, _ := filepath.Abs(sourcePath)
	for _, file := range foundFiles {
		if filepath.Ext(file.Name()) == allowedExtension {
			filePath := filepath.Join(path, file.Name())
			filePaths = append(filePaths, filePath)
			fmt.Println(filePath)
		}

	}
	return filePaths
}
