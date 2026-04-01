package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsWalker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "moti-test-walker-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	walker := &FsWalker{}
	var found bool
	err = walker.WalkDir(tmpDir, func(path string, err error) error {
		if path == "test.txt" {
			found = true
		}
		return nil
	})
	require.NoError(t, err)
	assert.True(t, found)
}
