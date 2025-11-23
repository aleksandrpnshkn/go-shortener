package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

func TestGetURLByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	existedCode := types.Code("tEsT")
	fullURL := types.OriginalURL("http://example.com")

	user := users.User{
		ID: 123,
	}

	urlsStorage := urls.NewMemoryStorage()
	urlsStorage.Set(context.Background(), urls.ShortenedURL{
		Code:        existedCode,
		OriginalURL: fullURL,
	}, user.ID)

	followedPublisher := audit.NewPublisher([]audit.Observer{})
	unshortener := services.NewUnshortener(urlsStorage, followedPublisher)

	t.Run("existed short url", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().FromUserContext(gomock.Any()).Return(user.ID, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+string(existedCode), nil)
		req.SetPathValue("code", string(existedCode))

		GetURLByCode(auther, unshortener)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "has redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Equal(t, string(fullURL), res.Header.Get("Location"), "redirects to original url")
	})

	t.Run("unknown short url", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().FromUserContext(gomock.Any()).Return(user.ID, nil)

		unknownCode := "uNkNoWn"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+unknownCode, nil)
		req.SetPathValue("code", unknownCode)

		GetURLByCode(auther, unshortener)(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "has no redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Empty(t, res.Header.Get("Location"))
	})
}
