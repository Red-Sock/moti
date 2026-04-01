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

func TestGitRepo_ReadFile(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	t.Run("success", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		revision := models.Revision{CommitHash: "hash1"}
		fileName := "buf.work.yaml"

		consoleMock.RunCmdMock.Inspect(func(ctx context.Context, dir string, command string, commandParams ...string) {
			assert.Equal(t, "/tmp", dir)
			assert.Equal(t, "git", command)
			assert.Equal(t, []string{"cat-file", "-p", "hash1:buf.work.yaml"}, commandParams)
		}).Return("file content", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		content, err := repo.ReadFile(ctx, revision, fileName)
		require.NoError(t, err)
		assert.Equal(t, "file content", content)
	})

	t.Run("file not found", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		revision := models.Revision{CommitHash: "hash1"}
		fileName := "nonexistent.file"

		consoleMock.RunCmdMock.Return("", assert.AnError)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		_, err := repo.ReadFile(ctx, revision, fileName)
		require.ErrorIs(t, err, models.ErrFileNotFound)
	})
}
