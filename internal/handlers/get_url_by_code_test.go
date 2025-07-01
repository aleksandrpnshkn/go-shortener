package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLByCode(t *testing.T) {
	existedCode := "tEsT"
	fullURL := "http://example.com"

	fullURLsStorage := services.NewFullURLsStorage()
	fullURLsStorage.Set(services.Code(existedCode), services.FullURL(fullURL))

	t.Run("existed short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+existedCode, nil)
		req.SetPathValue("code", existedCode)

		GetURLByCode(fullURLsStorage)(w, req)

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

		GetURLByCode(fullURLsStorage)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "has no redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Empty(t, res.Header.Get("Location"))
	})
}
