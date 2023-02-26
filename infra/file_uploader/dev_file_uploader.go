package file_uploader

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"ticken-event-service/models"
)

const uploadFilePath = ".uploads"

type DevFileUploader struct {
}

func NewDevFileUploader() (*DevFileUploader, error) {
	err := os.MkdirAll(uploadFilePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to init upload folder: %s", err.Error())
	}

	return &DevFileUploader{}, nil
}

// UploadFile is going to upload the given file to the AWS S3 bucket
func (uploader *DevFileUploader) UploadFile(file *models.File) (string, error) {
	newRandomFileName := uuid.New().String()

	relativePath := filepath.Join("uploads", newRandomFileName+file.GetExtension())

	newFile, err := os.Create(relativePath)
	if err != nil {
		return "", err
	}

	_, err = newFile.Write(*file.Content)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
