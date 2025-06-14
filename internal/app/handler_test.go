package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlByCode(t *testing.T) {
	existedCode := "tEsT"
	fullUrl := "http://example.com"

	app := application{
		config: config{
			hostname: "localhost",
			schema:   "http",
		},
		codesToURLs: map[string]string{existedCode: fullUrl},
	}

	t.Run("existed short url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+existedCode, nil)
		req.SetPathValue("code", existedCode)

		getUrlByCode(app)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "has redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Equal(t, fullUrl, res.Header.Get("Location"), "redirects to original url")
	})

	t.Run("unknown short url", func(t *testing.T) {
		unknownCode := "uNkNoWn"

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+unknownCode, nil)
		req.SetPathValue("code", unknownCode)

		getUrlByCode(app)(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "has no redirect")
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		assert.Empty(t, res.Header.Get("Location"))
	})
}

// Эндпоинт с методом POST и путём /
// возвращает ответ с кодом 201 и сокращённым URL как text/plain

// На любой некорректный запрос сервер должен возвращать ответ с кодом 400.
