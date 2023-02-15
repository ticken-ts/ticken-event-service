package file_uploader

type FileUploader struct {
}

func NewFileUploader() *FileUploader {
	return &FileUploader{}
}

// UploadFile is going to upload the given file to the AWS S3 bucket
func (uploader *FileUploader) UploadFile(file []byte, fileName string) (string, error) {
	panic("implement me")
}
