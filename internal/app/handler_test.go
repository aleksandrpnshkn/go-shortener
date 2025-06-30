package app

import (
	"net/http"
	"net/http/httptest"
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
