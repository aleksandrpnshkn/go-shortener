package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestDeleteUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := users.User{
		ID: 123,
	}

	auther := mocks.NewMockAuther(ctrl)
	auther.EXPECT().FromUserContext(gomock.Any()).Return(&user, nil)

	codeGenerator := services.NewTestGenerator("")

	urlsStorage := mocks.NewMockURLsStorage(ctrl)
	urlsStorage.EXPECT().DeleteManyByUserID(gomock.Any(), []types.Code{"foo", "bar"}, &user)

	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		"http://localhost",
	)

	t.Run("delete user urls success", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(`["foo", "bar"]`)
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", reqBody)

		var wg sync.WaitGroup
		DeleteUserURLs(shortener, zap.NewExample(), auther, &wg)(w, req)
		wg.Wait()

		res := w.Result()
		assert.Equal(t, http.StatusAccepted, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
