package sync

import (
	"io"
	"net/http"
)

func readBody(res *http.Response) string {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	return string(resBody)
}
