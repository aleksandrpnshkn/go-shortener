package audit

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoteLogger(t *testing.T) {
	t.Run("send audit entry", func(t *testing.T) {
		srv := httptest.NewServer(
			http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				rawLogs, err := io.ReadAll(req.Body)
				require.NoError(t, err, "should be no error while reading request")
				defer req.Body.Close()

				assert.Equal(
					t,
					`{"ts":12345678,"action":"shorten","user_id":"12315134","url":"https://mylongdomain.com/my/long/path/to/shorten/"}`,
					string(rawLogs),
					"check request content",
				)
			}),
		)
		defer srv.Close()

		remoteLogger := NewRemoteLogger(srv.Client(), srv.URL)
		err := remoteLogger.SendEntry(context.Background(), Entry{
			TimeTs:      time.Unix(12345678, 0).Unix(),
			Action:      "shorten",
			UserID:      "12315134",
			OriginalURL: "https://mylongdomain.com/my/long/path/to/shorten/",
		})

		assert.NoError(t, err, "should be no error")
	})
}
