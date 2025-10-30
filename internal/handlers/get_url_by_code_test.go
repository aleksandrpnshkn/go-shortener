package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLByCode(t *testing.T) {
	existedCode := types.Code("tEsT")
	fullURL := types.OriginalURL("http://example.com")

	user := users.User{
		ID: 123,
	}

	urlsStorage := urls.NewMemoryStorage()
	urlsStorage.Set(context.Background(), urls.ShortenedURL{
		Code:        existedCode,
		OriginalURL: fullURL,
	}, &user)

	t.Run("existed short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+string(existedCode), nil)
		req.SetPathValue("code", string(existedCode))

		GetURLByCode(urlsStorage)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "has redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Equal(t, string(fullURL), res.Header.Get("Location"), "redirects to original url")
	})

	t.Run("unknown short url", func(t *testing.T) {
		unknownCode := "uNkNoWn"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+unknownCode, nil)
		req.SetPathValue("code", unknownCode)

		GetURLByCode(urlsStorage)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "has no redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Empty(t, res.Header.Get("Location"))
	})
}
