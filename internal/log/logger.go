package log

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func NewLogger(level string) (*zap.Logger, error) {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zapLevel
	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
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
