package console

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBash_RunCmd(t *testing.T) {
	b := New()
	ctx := context.Background()
	wd, _ := os.Getwd()

	t.Run("success", func(t *testing.T) {
		out, err := b.RunCmd(ctx, wd, "echo", "hello")
		require.NoError(t, err)
		assert.Contains(t, out, "hello")
	})

	t.Run("error", func(t *testing.T) {
		_, err := b.RunCmd(ctx, wd, "nonexistentcommand")
		assert.Error(t, err)

		runErr, ok := err.(*RunError)
		require.True(t, ok)
		assert.Equal(t, "nonexistentcommand", runErr.Command)
		assert.NotEmpty(t, runErr.Stderr)
	})

	t.Run("with directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "moti-console-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		out, err := b.RunCmd(ctx, tmpDir, "pwd")
		require.NoError(t, err)
		// On macOS /var might be a symlink to /private/var
		assert.Contains(t, out, filepath.Base(tmpDir))
	})
}
