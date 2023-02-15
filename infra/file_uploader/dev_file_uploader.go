package file_uploader

import (
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"ticken-event-service/models"
)

type DevFileUploader struct {
}

func NewDevFileUploader() *DevFileUploader {
	return &DevFileUploader{}
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
