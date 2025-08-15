package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLPlain(t *testing.T) {
	fullURL := "http://example.com"
	code := "tEsT"

	codeGenerator := services.NewTestGenerator(code)
	URLsStorage := services.NewURLsTestStorage()
	shortener := services.NewShortener(
		codeGenerator,
		URLsStorage,
		"http://localhost",
	)

	t.Run("create short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := strings.NewReader(fullURL)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)

		CreateShortURLPlain(shortener)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = res.Body.Close()
		require.NoError(t, err)

		rawShortURL := string(resBody)
		assert.Equal(t, "http://localhost/tEsT", rawShortURL, "returned short url")

		storedURL, _ := URLsStorage.Get(context.Background(), services.Code(code))
		assert.Equal(t, fullURL, string(storedURL), "new code stored")
	})
}

func TestCreateShort(t *testing.T) {
	code := "tEsT"

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
			responseRawBody: `{"result":"http://localhost/tEsT"}`,
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
		codeGenerator := services.NewTestGenerator(code)
		URLsStorage := services.NewURLsTestStorage()
		shortener := services.NewShortener(
			codeGenerator,
			URLsStorage,
			"http://localhost",
		)

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody := strings.NewReader(test.requestRawBody)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)

			CreateShortURL(shortener)(w, req)

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
