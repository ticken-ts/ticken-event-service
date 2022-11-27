package utils

import (
	"archive/zip"
	"bytes"
)

func ZipFiles(filename2Content map[string][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for filename, content := range filename2Content {
		f, err := writer.Create(filename)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(content)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
