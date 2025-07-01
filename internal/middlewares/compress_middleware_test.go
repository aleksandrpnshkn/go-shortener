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

const testText = `{"hello": "world"}`
const testStatus = http.StatusOK

func testEchoHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		requestText, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		res.WriteHeader(testStatus)
		res.Write(requestText)
	})
}

func TestCompressMiddleware(t *testing.T) {
	handler := CompressMiddleware(testEchoHandler())

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
