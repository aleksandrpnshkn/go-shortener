package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type decompressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func (d *decompressReader) Read(p []byte) (n int, err error) {
	return d.gr.Read(p)
}

func (d *decompressReader) Close() error {
	if err := d.r.Close(); err != nil {
		return err
	}
	return d.gr.Close()
}

func newDecompressReader(req io.ReadCloser) (*decompressReader, error) {
	gzipReader, err := gzip.NewReader(req)
	if err != nil {
		return nil, err
	}

	decompressReader := &decompressReader{
		r:  req,
		gr: gzipReader,
	}

	return decompressReader, nil
}

func NewDecompressMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			contentEncoding := req.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")

			if sendsGzip {
				decompressedBody, err := newDecompressReader(req.Body)
				if err != nil {
					logger.Error("failed to decompress request", zap.Error(err))
					res.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer decompressedBody.Close()
				req.Body = decompressedBody
			}

			next.ServeHTTP(res, req)
		})
	}
}
