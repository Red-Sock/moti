package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/models"
)

type mockLockFile struct {
	readFunc func(moduleName string) (models.LockFileInfo, error)
}

func (m *mockLockFile) Read(moduleName string) (models.LockFileInfo, error) {
	return m.readFunc(moduleName)
}

func TestStorage_Paths(t *testing.T) {
	s := New("/tmp/moti", nil)

	t.Run("GetInstallDir", func(t *testing.T) {
		path := s.GetInstallDir("github.com/user/repo", "v1.0.0")
		assert.Equal(t, "/tmp/moti/mod/github.com/user/repo/v1.0.0", path)

		path = s.GetInstallDir("github.com/user/repo", "v1/2/3")
		assert.Equal(t, "/tmp/moti/mod/github.com/user/repo/v1-2-3", path)
	})

	t.Run("GetCacheDownloadPaths", func(t *testing.T) {
		module := models.Module{Name: "github.com/user/repo"}
		revision := models.Revision{Version: "v1.0.0"}
		paths := s.GetCacheDownloadPaths(module, revision)

		assert.Equal(t, "/tmp/moti/cache/download/github.com/user/repo", paths.CacheDownloadDir)
		assert.Equal(t, "/tmp/moti/cache/download/github.com/user/repo/v1.0.0.zip", paths.ArchiveFile)
	})

	t.Run("CreateCacheRepositoryDir", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "moti-storage-test-*")
		defer os.RemoveAll(tmpDir)
		s := New(tmpDir, nil)

		path, err := s.CreateCacheRepositoryDir("repo1")
		require.NoError(t, err)
		assert.Contains(t, path, tmpDir)

		_, err = os.Stat(path)
		assert.NoError(t, err)
	})

	t.Run("CreateCacheDownloadDir", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "moti-storage-test-*")
		defer os.RemoveAll(tmpDir)
		s := New(tmpDir, nil)

		paths := models.CacheDownloadPaths{
			CacheDownloadDir: filepath.Join(tmpDir, "download"),
		}
		err := s.CreateCacheDownloadDir(paths)
		require.NoError(t, err)

		_, err = os.Stat(paths.CacheDownloadDir)
		assert.NoError(t, err)
	})
}

func TestStorage_IsModuleInstalled(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "moti-storage-test-*")
	defer os.RemoveAll(tmpDir)

	t.Run("module not in lock file", func(t *testing.T) {
		mlf := &mockLockFile{
			readFunc: func(moduleName string) (models.LockFileInfo, error) {
				return models.LockFileInfo{}, models.ErrModuleNotFoundInLockFile
			},
		}
		s := New(tmpDir, mlf)
		installed, err := s.IsModuleInstalled(models.Module{Name: "m1"})
		require.NoError(t, err)
		assert.False(t, installed)
	})

	t.Run("module installed and matches", func(t *testing.T) {
		version := "v1.0.0"
		moduleName := "github.com/user/repo"
		installDir := filepath.Join(tmpDir, "mod", moduleName, version)
		err := os.MkdirAll(installDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(installDir, "file.proto"), []byte("syntax = \"proto3\";"), 0644)
		require.NoError(t, err)

		// Calculate actual hash
		s_temp := New(tmpDir, nil)
		hash, err := s_temp.GetInstalledModuleHash(moduleName, version)
		require.NoError(t, err)

		mlf := &mockLockFile{
			readFunc: func(name string) (models.LockFileInfo, error) {
				return models.LockFileInfo{
					Name:    moduleName,
					Version: version,
					Hash:    hash,
				}, nil
			},
		}
		s := New(tmpDir, mlf)
		installed, err := s.IsModuleInstalled(models.Module{Name: moduleName, Version: models.RequestedVersion(version)})
		require.NoError(t, err)
		assert.True(t, installed)
	})
}
