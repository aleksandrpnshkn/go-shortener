package audit

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileLogger(t *testing.T) {
	t.Run("write line", func(t *testing.T) {
		testFileName := "write_line_unit_test.log"
		fileLogger, err := NewFileLogger(testFileName)
		require.NoError(t, err, "file logger should be created")

		err = fileLogger.AddEntry(context.Background(), Entry{
			TimeTS:      time.Unix(12345678, 0).Unix(),
			Action:      "shorten",
			UserID:      "12315134",
			OriginalURL: "https://mylongdomain.com/my/long/path/to/shorten/",
		})
		require.NoError(t, err, "should add entry to log")

		err = fileLogger.Close()
		require.NoError(t, err, "should close logger")

		rawLogs, err := os.ReadFile(testFileName)
		require.NoError(t, err, "should read logs")
		assert.Equal(
			t,
			`{"ts":12345678,"action":"shorten","user_id":"12315134","url":"https://mylongdomain.com/my/long/path/to/shorten/"}
`,
			string(rawLogs),
			"check file content",
		)

		err = os.Remove(testFileName)
		require.NoError(t, err, "should clean remove file after test")
	})
}
