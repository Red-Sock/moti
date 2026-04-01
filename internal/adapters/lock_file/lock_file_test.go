package lockfile

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/models"
)

type mockDirWalker struct {
	files map[string]*bytes.Buffer
}

func (m *mockDirWalker) Open(name string) (io.ReadCloser, error) {
	buf, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func (m *mockDirWalker) Create(name string) (io.WriteCloser, error) {
	buf := &bytes.Buffer{}
	m.files[name] = buf
	return &nopWriteCloser{buf}, nil
}

func (m *mockDirWalker) WalkDir(callback func(path string, err error) error) error {
	for path := range m.files {
		if err := callback(path, nil); err != nil {
			return err
		}
	}
	return nil
}

type nopWriteCloser struct {
	*bytes.Buffer
}

func (n *nopWriteCloser) Close() error {
	return nil
}

func TestNewLockFile(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		dw := &mockDirWalker{
			files: map[string]*bytes.Buffer{
				"moti.lock": bytes.NewBufferString("module1 v1.0.0 hash1\nmodule2 v2.0.0 hash2\n"),
			},
		}

		lf, err := New(dw)
		require.NoError(t, err)
		assert.NotNil(t, lf)

		info, err := lf.Read("module1")
		require.NoError(t, err)
		assert.Equal(t, "module1", info.Name)
		assert.Equal(t, "v1.0.0", info.Version)
		assert.Equal(t, models.ModuleHash("hash1"), info.Hash)

		info, err = lf.Read("module2")
		require.NoError(t, err)
		assert.Equal(t, "v2.0.0", info.Version)

		_, err = lf.Read("nonexistent")
		assert.ErrorIs(t, err, models.ErrModuleNotFoundInLockFile)
	})

	t.Run("no lock file", func(t *testing.T) {
		dw := &mockDirWalker{
			files: make(map[string]*bytes.Buffer),
		}

		lf, err := New(dw)
		require.NoError(t, err)
		assert.NotNil(t, lf)

		_, err = lf.Read("any")
		assert.ErrorIs(t, err, models.ErrModuleNotFoundInLockFile)
	})
}

func TestLockFile_Write(t *testing.T) {
	dw := &mockDirWalker{
		files: make(map[string]*bytes.Buffer),
	}

	lf, err := New(dw)
	require.NoError(t, err)

	err = lf.Write("module-b", "v2", "hash-b")
	require.NoError(t, err)

	err = lf.Write("module-a", "v1", "hash-a")
	require.NoError(t, err)

	// Check if file was written correctly (and sorted)
	buf, ok := dw.files["moti.lock"]
	require.True(t, ok)
	
	expected := "module-a v1 hash-a\nmodule-b v2 hash-b\n"
	assert.Equal(t, expected, buf.String())

	// Check if cache was updated
	info, err := lf.Read("module-a")
	require.NoError(t, err)
	assert.Equal(t, "v1", info.Version)
}
