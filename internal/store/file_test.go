package store

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {
	_, testFileName, _, _ := runtime.Caller(0)
	testFilesDirPath := filepath.Dir(testFileName) + "/test_files"

	t.Run("entries loaded", func(t *testing.T) {
		fileStorage, err := NewFileStorage(testFilesDirPath + "/test_load.txt")
		require.NoError(t, err)
		defer fileStorage.Close()

		assert.Equal(t, 2, fileStorage.lastID, "last id loaded")

		originalURL, isFound := fileStorage.Get("test2")

		assert.True(t, isFound, "entry found")
		assert.Equal(t, "http://example2.com", originalURL, "original url loaded")
	})
}
