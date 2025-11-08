package handlers

import "net/http"

// FallbackHandler - хендлер для неизвестных маршрутов.
func FallbackHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain")
		res.WriteHeader(http.StatusBadRequest)
	}
}
