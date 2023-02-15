package file_uploader

import (
	"ticken-event-service/models"
)

type FileUploader struct {
}

func NewFileUploader() *FileUploader {
	return &FileUploader{}
}

// UploadFile is going to upload the given file to the AWS S3 bucket
func (uploader *FileUploader) UploadFile(file *models.File) (string, error) {
	panic("implement me")
}
