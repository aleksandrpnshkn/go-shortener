package compress

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressMiddleware(t *testing.T) {
	testText := `{"hello": "world"}`
	testStatus := http.StatusOK
	handler := DecompressMiddleware(testEchoHandler(testStatus))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("client sent gzip", func(t *testing.T) {
		gzippedBuf, err := gzipText(testText)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", srv.URL, gzippedBuf)
		req.RequestURI = ""
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, testStatus, resp.StatusCode, "successful status")
		defer resp.Body.Close()

		rawResponseText, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.JSONEq(t, testText, string(rawResponseText), "server echoed the same json")
	})
}
