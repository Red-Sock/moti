package generate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/mocks"
	"go.redsock.ru/moti/internal/models"
)

func TestGenerate(t *testing.T) {
	ctx := t.Context()

	mConsole := mocks.NewConsoleMock(t)

	expectedParams := []string{
		"-I test",
		"test/test/file.proto",
	}

	mConsole.RunCmdMock.Expect(ctx, ".", protocBin, expectedParams...).
		Return("", nil)

	mStorage := mocks.NewIStorageMock(t)

	mStorage.IsModuleInstalledMock.Return(true, nil)
	mStorage.GetInstallDirMock.Set(
		func(moduleName string, revisionVersion string) string {
			return "/tmp/mod/" + moduleName
		})

	mLockFile := mocks.NewILockFileMock(t)
	mLockFile.ReadMock.Return(models.LockFileInfo{Version: "v1.0.0"}, nil)

	mWalker := mocks.NewIWalkerMock(t)
	mWalker.WalkDirMock.Set(
		func(root string, callback func(path string, err error) error) error {
			return callback("test/file.proto", nil)
		})

	c := &Core{
		Env: commands.Env{
			MotiConfig: config.Config{
				Generate: []config.Generate{
					{
						Inputs: []config.Input{
							{Directory: "test"},
						},
					},
				},
			},
			Console:  mConsole,
			Storage:  mStorage,
			LockFile: mLockFile,
			WorkDir:  ".",
		},
		Walker: mWalker,
	}

	err := c.Generate(ctx)
	require.NoError(t, err)
}

func TestGenerate_MultipleInputs(t *testing.T) {
	ctx := t.Context()

	expectedParams1 := []string{
		"-I test1",
		"test1/test1/file.proto",
	}

	expectedParams2 := []string{
		"-I test2",
		"test2/test2/file.proto",
	}

	mConsole := mocks.NewConsoleMock(t)
	mConsole.
		RunCmdMock.
		When(ctx, ".", protocBin, expectedParams1...).
		Then("", nil).
		RunCmdMock.
		When(ctx, ".", protocBin, expectedParams2...).
		Then("", nil)

	mStorage := mocks.NewIStorageMock(t)
	mLockFile := mocks.NewILockFileMock(t)
	mWalker := mocks.NewIWalkerMock(t)
	mWalker.WalkDirMock.Set(func(root string, callback func(path string, err error) error) error {
		return callback(root+"/file.proto", nil)
	})

	c := &Core{
		Env: commands.Env{
			MotiConfig: config.Config{
				Generate: []config.Generate{
					{
						Inputs: []config.Input{
							{Directory: "test1"},
							{Directory: "test2"},
						},
					},
				},
			},
			Console:  mConsole,
			Storage:  mStorage,
			LockFile: mLockFile,
			WorkDir:  ".",
		},
		Walker: mWalker,
	}

	err := c.Generate(ctx)
	require.NoError(t, err)
}
