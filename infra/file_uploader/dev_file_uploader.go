package file_uploader

type DevFileUploader struct {
}

func NewDevFileUploader() *DevFileUploader {
	return &DevFileUploader{}
}

// UploadFile is going to upload the given file to the AWS S3 bucket
func (uploader *DevFileUploader) UploadFile(file []byte, fileName string) (string, error) {
	panic("implement me")
}
