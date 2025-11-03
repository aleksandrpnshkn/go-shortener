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
	"github.com/aleksandrpnshkn/go-shortener/internal/services/audit"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateShortBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := users.User{
		ID: 123,
	}

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
		codesReserver := mocks.NewMockCodesReserver(ctrl)
		codesReserver.EXPECT().GetCode(gomock.Any()).AnyTimes().Return(types.Code("tEsT1"), nil)

		urlsStorage := urls.NewMemoryStorage()
		shortener := services.NewShortener(
			context.Background(),
			codesReserver,
			urlsStorage,
			"http://localhost",
			audit.NewPublisher([]audit.Observer{}),
		)

		t.Run(test.testName, func(t *testing.T) {
			auther := mocks.NewMockAuther(ctrl)
			auther.EXPECT().FromUserContext(gomock.Any()).Return(&user, nil)

			w := httptest.NewRecorder()
			reqBody := strings.NewReader(test.requestRawBody)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", reqBody)

			CreateShortURLBatch(shortener, zap.NewExample(), auther)(w, req)

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
