package models

import "mime"

type File struct {
	Content  *[]byte
	MimeType string
}

func NewFile(content *[]byte, mime string) *File {
	return &File{
		Content:  content,
		MimeType: mime,
	}
}

// GetExtension returns the extension of the file based on the mime type, includes the dot
// Example: "image/png" -> ".png"
func (file *File) GetExtension() string {
	println("getting extension of: " + file.MimeType)
	extensions, err := mime.ExtensionsByType(file.MimeType)
	if err != nil || extensions == nil {
		return ""
	}
	println("extension: " + extensions[0])
	return extensions[0]
}
