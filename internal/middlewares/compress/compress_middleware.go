// Модуль compress содержит в себе middleware для распаковки запроса/сжатия ответов сервера.
// Поддерживает формат gzip.
package compress

import (
	"compress/gzip"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type compressWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(b []byte) (int, error) {
	return c.gw.Write(b)
}

func (c *compressWriter) Close() error {
	return c.gw.Close()
}

func newCompressWriter(res http.ResponseWriter) (*compressWriter, error) {
	gzipWriter, err := gzip.NewWriterLevel(res, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}

	compressWriter := &compressWriter{
		w:  res,
		gw: gzipWriter,
	}

	return compressWriter, nil
}

func NewCompressMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			acceptEncoding := req.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")

			contentType := req.Header.Get("Content-Type")
			isGzipableContentType := strings.Contains(contentType, "application/json") ||
				strings.Contains(contentType, "text/html")

			if supportsGzip && isGzipableContentType {
				compressWriter, err := newCompressWriter(res)
				if err != nil {
					logger.Error("failed to compress response", zap.Error(err))
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				defer compressWriter.Close()
				res = compressWriter

				res.Header().Set("Content-Encoding", "gzip")
			}

			next.ServeHTTP(res, req)
		})
	}
}
