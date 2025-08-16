package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateShortURLPlain(t *testing.T) {
	fullURL := "http://example.com"

	codeGenerator := services.NewTestGenerator("tEsT")
	urlsStorage := store.NewMemoryStorage()
	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		"http://localhost",
	)

	t.Run("create short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(fullURL)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)

		CreateShortURLPlain(shortener, zap.NewExample())(w, req)

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
	codePrefix := "tEsT"

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
		urlsStorage := store.NewMemoryStorage()
		shortener := services.NewShortener(
			codeGenerator,
			urlsStorage,
			"http://localhost",
		)

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody := strings.NewReader(test.requestRawBody)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)

			CreateShortURL(shortener, zap.NewExample())(w, req)

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
	codePrefix := "tEsT"
	originalURL := "http://example.com"

	codeGenerator := services.NewTestGenerator(codePrefix)
	urlsStorage := store.NewMemoryStorage()
	urlsStorage.Set(context.Background(), store.ShortenedURL{
		Code:        "test123",
		OriginalURL: originalURL,
	})
	shortener := services.NewShortener(
		codeGenerator,
		urlsStorage,
		"http://localhost",
	)

	t.Run("create duplicate url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(`{"url":"` + originalURL + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)

		CreateShortURL(shortener, zap.NewExample())(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusConflict, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
