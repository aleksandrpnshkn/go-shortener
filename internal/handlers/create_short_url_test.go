package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURL(t *testing.T) {
	fullURL := "http://example.com"

	codeGenerator := services.NewCodeGenerator(3)
	fullURLsStorage := services.NewFullURLsStorage()
	shortener := services.NewShortener(
		*codeGenerator,
		*fullURLsStorage,
		"http://localhost",
	)

	t.Run("create short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(fullURL)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)

		CreateShortURL(shortener)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)

		rawShortURL := string(resBody)
		assert.Contains(t, rawShortURL, "http://localhost", "contains hostname")

		shortURL, err := url.Parse(rawShortURL)
		require.NoError(t, err, "response contains correct short url")
		code := strings.TrimLeft(shortURL.Path, "/")
		assert.NotEmpty(t, code, "path has code")

		storedURL, _ := fullURLsStorage.Get(services.Code(code))
		assert.Equal(t, fullURL, string(storedURL), "new code stored")
	})
}
