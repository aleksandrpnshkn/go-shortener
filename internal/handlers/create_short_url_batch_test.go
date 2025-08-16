package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortBatch(t *testing.T) {
	codePrefix := "tEsT"

	tests := []struct {
		testName        string
		statusCode      int
		requestRawBody  string
		responseRawBody string
	}{
		{
			testName:        "create one short url in batch",
			statusCode:      http.StatusCreated,
			requestRawBody:  `[{"correlation_id": "c1", "original_url":"http://example.com"}]`,
			responseRawBody: `[{"correlation_id": "c1", "short_url":"http://localhost/tEsT1"}]`,
		},
		{
			testName:        "invalid batch json",
			statusCode:      http.StatusBadRequest,
			requestRawBody:  `{"foo":"bar"}`,
			responseRawBody: `{"error":{"message":"bad request"}}`,
		},
		{
			testName:        "malformatted batch json",
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
			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", reqBody)

			CreateShortURLBatch(shortener, logs.NewTestLogger())(w, req)

			res := w.Result()
			assert.Equal(t, test.statusCode, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			assert.JSONEq(t, test.responseRawBody, string(resBody), "response json")
		})
	}
}
