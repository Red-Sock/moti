package generate

import (
	"context"
	"strings"
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
	mConsole.RunCmdMock.Set(func(ctx context.Context, dir string, command string, commandParams ...string) (string, error) {
		require.Equal(t, "protoc", command)
		// Verify some args
		foundI := false
		foundProto := false
		for _, arg := range commandParams {
			if arg == "-I ." {
				foundI = true
			}
			if strings.Contains(arg, "test/*.proto") {
				foundProto = true
			}
		}
		require.True(t, foundI)
		require.True(t, foundProto)
		return "", nil
	})

	mStorage := mocks.NewIStorageMock(t)
	mStorage.IsModuleInstalledMock.Return(true, nil)
	mStorage.GetInstallDirMock.Set(func(moduleName string, revisionVersion string) string {
		return "/tmp/mod/" + moduleName
	})

	mLockFile := mocks.NewILockFileMock(t)
	mLockFile.ReadMock.Return(models.LockFileInfo{Version: "v1.0.0"}, nil)

	mWalker := mocks.NewIWalkerMock(t)
	mWalker.WalkDirMock.Set(func(root string, path string, callback func(path string, err error) error) error {
		return callback("test/file.proto", nil)
	})

	c := &Core{
		Env: commands.Env{
			MotiConfig: config.Config{
				Generate: config.Generate{
					Inputs: []config.Input{
						{Directory: "test"},
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
