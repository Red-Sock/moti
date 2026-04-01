package fs

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirWalker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "moti-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)

	dw := NewFSWalker(tmpDir, ".")

	t.Run("Open", func(t *testing.T) {
		rc, err := dw.Open("file1.txt")
		require.NoError(t, err)
		defer rc.Close()

		content, err := io.ReadAll(rc)
		require.NoError(t, err)
		assert.Equal(t, "content1", string(content))
	})

	t.Run("Create", func(t *testing.T) {
		wc, err := dw.Create("newfile.txt")
		require.NoError(t, err)

		_, err = wc.Write([]byte("newcontent"))
		require.NoError(t, err)
		err = wc.Close()
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "newfile.txt"))
		require.NoError(t, err)
		assert.Equal(t, "newcontent", string(content))
	})

	t.Run("WalkDir", func(t *testing.T) {
		var paths []string
		err := dw.WalkDir(func(path string, err error) error {
			require.NoError(t, err)
			paths = append(paths, path)
			return nil
		})
		require.NoError(t, err)

		assert.Contains(t, paths, ".")
		assert.Contains(t, paths, "file1.txt")
		assert.Contains(t, paths, "subdir")
		assert.Contains(t, paths, "subdir/file2.txt")
	})
}
