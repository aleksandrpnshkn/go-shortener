package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressMiddleware(t *testing.T) {
	testText := `{"hello": "world"}`
	testStatus := http.StatusOK
	handler := CompressMiddleware(testEchoHandler(testStatus))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("server compressed json", func(t *testing.T) {
		buf := bytes.NewBufferString(testText)

		req := httptest.NewRequest("POST", srv.URL, buf)
		req.RequestURI = ""
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, testStatus, resp.StatusCode, "successful status")
		defer resp.Body.Close()

		gzipReader, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		rawResponseText, err := io.ReadAll(gzipReader)
		require.NoError(t, err)
		assert.Equal(t, testText, string(rawResponseText), "server echoed the same json")
	})

	t.Run("server not compressed plain text", func(t *testing.T) {
		buf := bytes.NewBufferString(testText)

		req := httptest.NewRequest("POST", srv.URL, buf)
		req.RequestURI = ""
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, testStatus, resp.StatusCode, "successful status")
		defer resp.Body.Close()

		rawResponseText, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, testText, string(rawResponseText), "server echoed the same json")
	})
}
