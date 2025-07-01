package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
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

type compressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.gr.Close()
}

func newCompressReader(req io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(req)
	if err != nil {
		return nil, err
	}

	compressWriter := &compressReader{
		r:  req,
		gr: gzipReader,
	}

	return compressWriter, nil
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		contentType := req.Header.Get("Content-Type")
		isGzipableContentType := strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "text/html")

		if supportsGzip && isGzipableContentType {
			compressRes, err := newCompressWriter(res)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer compressRes.Close()
			res = compressRes

			res.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			compressBody, err := newCompressReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer compressBody.Close()
			req.Body = compressBody
		}

		next.ServeHTTP(res, req)
	})
}
