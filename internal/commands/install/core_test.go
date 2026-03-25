package install

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/mocks"
	"go.redsock.ru/moti/internal/models"
)

func TestInstall(t *testing.T) {
	ctx := t.Context()

	mStorage := mocks.NewIStorageMock(t)
	mStorage.IsModuleInstalledMock.Return(false, nil)
	mStorage.CreateCacheRepositoryDirMock.Set(func(name string) (string, error) {
		return "/tmp/cache/" + name, nil
	})
	mStorage.GetCacheDownloadPathsMock.Set(func(module models.Module, revision models.Revision) models.CacheDownloadPaths {
		return models.CacheDownloadPaths{
			CacheDownloadDir: "/tmp/download/" + module.Name,
		}
	})
	mStorage.CreateCacheDownloadDirMock.Return(nil)
	mStorage.InstallMock.Return("some-hash", nil)

	mModuleConfig := mocks.NewIModuleConfigMock(t)
	mModuleConfig.ReadFromRepoMock.Return(models.ModuleConfig{}, nil)

	mLockFile := mocks.NewILockFileMock(t)
	mLockFile.WriteMock.Set(func(moduleName string, revisionVersion string, hash models.ModuleHash) error {
		require.Equal(t, "github.com/test/repo", moduleName)
		require.Equal(t, "v1.0.0", revisionVersion)
		require.Equal(t, models.ModuleHash("some-hash"), hash)
		return nil
	})

	mRepo := mocks.NewRepoMock(t)
	mRepo.ReadRevisionMock.Return(models.Revision{Version: "v1.0.0"}, nil)
	mRepo.FetchMock.Return(nil)
	mRepo.ArchiveMock.Return(nil)

	mRepoFactory := mocks.NewRepoFactoryMock(t)
	mRepoFactory.NewMock.Return(mRepo, nil)

	c := &Core{
		Env: commands.Env{
			MotiConfig: config.Config{
				Deps: []string{"github.com/test/repo@v1.0.0"},
			},
			Storage:      mStorage,
			ModuleConfig: mModuleConfig,
			LockFile:     mLockFile,
		},
		RepoFactory: mRepoFactory,
	}

	err := c.Install(ctx)
	require.NoError(t, err)
}
