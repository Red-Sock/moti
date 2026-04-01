package install

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/mocks"
	"go.redsock.ru/moti/internal/models"
)

func TestInstall_Binaries(t *testing.T) {
	ctx := t.Context()

	mConsole := mocks.NewConsoleMock(t)
	mConsole.RunCmdMock.When(ctx, "/tmp", "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11").Then("installed", nil)
	mConsole.RunCmdMock.When(ctx, "/tmp", "protoc-gen-go --version").Then("v1.36.11", nil)

	c := &Core{
		Env: commands.Env{
			WorkDir: "/tmp",
			MotiConfig: config.Config{
				Binaries: config.Binaries{
					Install: []struct {
						Go config.GoBin `json:"go" yaml:"go"`
					}{
						{
							Go: config.GoBin{
								Module:           "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11",
								VersionCheckArgs: "--version",
							},
						},
					},
				},
			},
			Console: mConsole,
		},
	}

	err := c.Install(ctx)
	require.NoError(t, err)
}

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

	mConsole := mocks.NewConsoleMock(t)

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
			Console:      mConsole,
		},
		RepoFactory: mRepoFactory,
	}

	err := c.Install(ctx)
	require.NoError(t, err)
}
