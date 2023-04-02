package file

import "mime"

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
