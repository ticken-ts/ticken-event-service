package file_uploader

import (
	"ticken-event-service/utils/file"
)

type FileUploader struct {
}

func NewFileUploader() (*FileUploader, error) {
	return &FileUploader{}, nil
}

// UploadFile is going to upload the given file to the AWS S3 bucket
func (uploader *FileUploader) UploadFile(file *file.File) (string, error) {
	panic("implement me")
}
