package file

import (
	"fmt"
	"io"
	"mime"
	"net/http"
)

type File struct {
	Filename string
	Content  []byte
	MimeType string
}

// GetExtension returns the extension of the
// file based on the mime type.
// Note: the extension includes the dot
// Example: "image/png" -> ".png"
func (file *File) GetExtension() string {
	extensions, err := mime.ExtensionsByType(file.MimeType)
	if err != nil || extensions == nil {
		return ""
	}
	return extensions[0]
}

func Download(URL string) (*File, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 response code")
	}

	imageContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	contentDisposition := response.Header.Get("Content-Disposition")

	var filename = ""
	if _, params, err := mime.ParseMediaType(contentDisposition); err == nil {
		filename = params["filename"]
	}

	return &File{
		Content:  imageContent,
		Filename: filename,
		MimeType: response.Header.Get("Content-Type"),
	}, nil
}
