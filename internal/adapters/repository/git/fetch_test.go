package git_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/mocks"
	"go.redsock.ru/moti/internal/models"
)

func TestGitRepo_Fetch(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	t.Run("success", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		revision := models.Revision{CommitHash: "hash1"}

		consoleMock.RunCmdMock.Inspect(func(ctx context.Context, dir string, command string, commandParams ...string) {
			assert.Equal(t, "/tmp", dir)
			assert.Equal(t, "git", command)
			assert.Equal(t, []string{"fetch", "-f", "origin", "--depth=1", "hash1"}, commandParams)
		}).Return("", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		err := repo.Fetch(ctx, revision)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		revision := models.Revision{CommitHash: "hash1"}

		consoleMock.RunCmdMock.Return("", assert.AnError)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		err := repo.Fetch(ctx, revision)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "adapters.RunCmd (fetch)")
	})
}
