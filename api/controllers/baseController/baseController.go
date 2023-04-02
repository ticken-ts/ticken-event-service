package baseController

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"ticken-event-service/services"
	"ticken-event-service/utils/file"
)

type BaseController struct {
	Validator       *validator.Validate
	ServiceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *BaseController {
	return &BaseController{
		Validator:       validator.New(),
		ServiceProvider: serviceProvider,
	}
}

func (controller *BaseController) ReadAsset(fileHeader *multipart.FileHeader) (*file.File, error) {
	fileContent, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileHeader.Size)

	if _, err = fileContent.Read(bytes); err != nil {
		return nil, err
	}

	return &file.File{
		Content:  bytes,
		Filename: fileHeader.Filename,
		MimeType: fileHeader.Header.Get("Content-Type"),
	}, nil
}
