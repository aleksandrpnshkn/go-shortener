package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLByCode(t *testing.T) {
	existedCode := "tEsT"
	fullURL := "http://example.com"

	app := application{
		config: &config.Config{
			ServerAddr:    "localhost",
			PublicBaseURL: "http://localhost",
		},
		codesToURLs: map[string]string{existedCode: fullURL},
	}

	t.Run("existed short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+existedCode, nil)
		req.SetPathValue("code", existedCode)

		getURLByCode(app)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "has redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Equal(t, fullURL, res.Header.Get("Location"), "redirects to original url")
	})

	t.Run("unknown short url", func(t *testing.T) {
		unknownCode := "uNkNoWn"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+unknownCode, nil)
		req.SetPathValue("code", unknownCode)

		getURLByCode(app)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "has no redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Empty(t, res.Header.Get("Location"))
	})
}

func TestCreateShortURL(t *testing.T) {
	fullURL := "http://example.com"

	app := application{
		config: &config.Config{
			ServerAddr:    "localhost",
			PublicBaseURL: "http://localhost",
		},
		codesToURLs: map[string]string{},
	}

	t.Run("create short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(fullURL)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)

		createShortURL(app)(w, req)

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
		assert.NotEmpty(t, code, "has code")
		assert.Equal(t, 1, len(app.codesToURLs), "new code stored")
		assert.Contains(t, app.codesToURLs, code, "new code stored")
	})
}

func TestGetFallback(t *testing.T) {
	t.Run("existed short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)

		fallbackHandler()(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	})
}
