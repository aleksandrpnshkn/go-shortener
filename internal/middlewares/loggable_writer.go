package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggableResponseWriter struct {
		writer       http.ResponseWriter
		responseData *responseData
	}
)

func (w *loggableResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *loggableResponseWriter) Write(b []byte) (int, error) {
	size, err := w.writer.Write(b)
	w.responseData.size += size
	return size, err
}

func (w *loggableResponseWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
	w.responseData.status = statusCode
}

func newLoggableResponseWriter(w *http.ResponseWriter) *loggableResponseWriter {
	responseData := responseData{}

	writer := loggableResponseWriter{
		writer:       *w,
		responseData: &responseData,
	}

	return &writer
}

func NewRequestMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()

			loggableResponseWriter := newLoggableResponseWriter(&res)

			next.ServeHTTP(loggableResponseWriter, req)

			statusCode := loggableResponseWriter.responseData.status
			responseSize := loggableResponseWriter.responseData.size
			duration := time.Since(start)

			logger.Info(req.Method+" "+req.RequestURI,
				zap.Int64("duration", duration.Microseconds()),
				zap.Int("status", statusCode),
				zap.Int("size", responseSize),
			)
		})
	}
}
