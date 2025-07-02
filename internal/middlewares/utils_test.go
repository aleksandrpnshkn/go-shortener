package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func gzipText(text string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(buf)
	_, err := gzipWriter.Write([]byte(text))
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	return buf, err
}

func testEchoHandler(statusCode int) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		requestText, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		res.WriteHeader(statusCode)
		res.Write(requestText)
	})
}
