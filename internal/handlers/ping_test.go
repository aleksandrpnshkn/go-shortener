package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
)

func TestPingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("test successful ping", func(t *testing.T) {
		urlsStorage := mocks.NewMockURLsStorage(ctrl)
		urlsStorage.EXPECT().Ping(context.Background()).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)

		PingHandler(context.Background(), urlsStorage, zap.NewExample())(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode, "status successful")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	})

	t.Run("test failed ping", func(t *testing.T) {
		urlsStorage := mocks.NewMockURLsStorage(ctrl)
		urlsStorage.EXPECT().Ping(context.Background()).Return(errors.New("something wrong"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)

		PingHandler(context.Background(), urlsStorage, zap.NewExample())(w, req)

		res := w.Result()
		err := res.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "server error occurred")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	})
}
