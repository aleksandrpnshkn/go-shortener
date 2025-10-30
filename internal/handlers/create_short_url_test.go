package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateShortURLPlain(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fullURL := "http://example.com"
	user := users.User{
		ID: 123,
	}

	auther := mocks.NewMockAuther(ctrl)
	auther.EXPECT().FromUserContext(gomock.Any()).Return(&user, nil)

	codeGenerator := services.NewTestGenerator("tEsT")
	urlsStorage := urls.NewMemoryStorage()
	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		"http://localhost",
	)

	t.Run("create short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(fullURL)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)

		CreateShortURLPlain(shortener, zap.NewExample(), auther)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)

		rawShortURL := string(resBody)
		assert.Equal(t, "http://localhost/tEsT1", rawShortURL, "returned short url")

		storedURL, _ := urlsStorage.Get(context.Background(), "tEsT1")
		assert.Equal(t, fullURL, string(storedURL), "new code stored")
	})
}

func TestCreateShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	codePrefix := "tEsT"

	user := users.User{
		ID: 123,
	}

	auther := mocks.NewMockAuther(ctrl)
	auther.EXPECT().FromUserContext(gomock.Any()).Return(&user, nil)

	tests := []struct {
		testName        string
		statusCode      int
		requestRawBody  string
		responseRawBody string
	}{
		{
			testName:        "create short url",
			statusCode:      http.StatusCreated,
			requestRawBody:  `{"url":"http://example.com"}`,
			responseRawBody: `{"result":"http://localhost/tEsT1"}`,
		},
		{
			testName:        "invalid json",
			statusCode:      http.StatusBadRequest,
			requestRawBody:  `{"foo":"bar"}`,
			responseRawBody: `{"error":{"message":"bad request"}}`,
		},
		{
			testName:        "malformatted json",
			statusCode:      http.StatusBadRequest,
			requestRawBody:  "}}}}}",
			responseRawBody: `{"error":{"message":"bad request"}}`,
		},
	}

	for _, test := range tests {
		codeGenerator := services.NewTestGenerator(codePrefix)
		urlsStorage := urls.NewMemoryStorage()
		shortener := services.NewShortener(
			codeGenerator,
			urlsStorage,
			"http://localhost",
		)

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody := strings.NewReader(test.requestRawBody)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)

			CreateShortURL(shortener, zap.NewExample(), auther)(w, req)

			res := w.Result()
			assert.Equal(t, test.statusCode, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			assert.JSONEq(t, string(resBody), test.responseRawBody, "response json")
		})
	}
}

func TestCreateShortDuplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	codePrefix := "tEsT"
	originalURL := "http://example.com"

	user := users.User{
		ID: 123,
	}

	auther := mocks.NewMockAuther(ctrl)
	auther.EXPECT().FromUserContext(gomock.Any()).Return(&user, nil)

	codeGenerator := services.NewTestGenerator(codePrefix)
	urlsStorage := urls.NewMemoryStorage()
	urlsStorage.Set(context.Background(), urls.ShortenedURL{
		Code:        "test123",
		OriginalURL: originalURL,
	}, &user)
	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		"http://localhost",
	)

	t.Run("create duplicate url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(`{"url":"` + originalURL + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)

		CreateShortURL(shortener, zap.NewExample(), auther)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusConflict, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
