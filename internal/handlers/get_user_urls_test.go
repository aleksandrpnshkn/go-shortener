package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
)

func TestGetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := users.User{
		ID: 123,
	}

	t.Run("user has some urls", func(t *testing.T) {
		testURL := urls.ShortenedURL{
			Code:        "tEsT",
			OriginalURL: "http://example.com",
		}
		userURLs := []urls.ShortenedURL{testURL}

		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().FromUserContext(gomock.Any()).Return(user.ID, nil)

		urlsStorage := mocks.NewMockURLsStorage(ctrl)
		urlsStorage.EXPECT().GetByUserID(gomock.Any(), user.ID).Return(userURLs, nil)
		shortener := services.NewShortener(
			context.Background(),
			mocks.NewMockCodesReserver(ctrl),
			urlsStorage,
			"http://localhost",
			audit.NewPublisher([]audit.Observer{}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		GetUserURLs(shortener, zap.NewExample(), auther)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)

		assert.JSONEq(t, `[{
			"original_url": "http://example.com", 
			"short_url": "http://localhost/tEsT"
		}]`, string(resBody))
	})

	t.Run("user has no urls", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().FromUserContext(gomock.Any()).Return(user.ID, nil)

		urlsStorage := mocks.NewMockURLsStorage(ctrl)
		urlsStorage.EXPECT().GetByUserID(gomock.Any(), user.ID).Return([]urls.ShortenedURL{}, nil)

		shortener := services.NewShortener(
			context.Background(),
			mocks.NewMockCodesReserver(ctrl),
			urlsStorage,
			"http://localhost",
			audit.NewPublisher([]audit.Observer{}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		GetUserURLs(shortener, zap.NewExample(), auther)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)

		assert.Empty(t, resBody)
	})
}
