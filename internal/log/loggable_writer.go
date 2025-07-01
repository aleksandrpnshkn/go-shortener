package log

import "net/http"

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
