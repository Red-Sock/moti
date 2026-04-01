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

func TestGitRepo_ReadRevision(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	t.Run("read by tag", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		consoleMock.RunCmdMock.Inspect(func(ctx context.Context, dir string, command string, commandParams ...string) {
			assert.Equal(t, "git", command)
			assert.Equal(t, "ls-remote", commandParams[0])
			assert.Equal(t, "origin", commandParams[1])
			assert.Equal(t, "v1.0.0", commandParams[2])
		}).Return("hash1\trefs/tags/v1.0.0\n", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		rev, err := repo.ReadRevision(ctx, "v1.0.0")
		require.NoError(t, err)
		assert.Equal(t, "hash1", rev.CommitHash)
		assert.Equal(t, "v1.0.0", rev.Version)
	})

	t.Run("read omitted (latest) with tag", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		// First call to get HEAD
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "ls-remote", "origin", "HEAD").
			Then("latest-hash\tHEAD\n", nil)
		// Second call to get tags for the commit
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "ls-remote", "origin").
			Then("latest-hash\trefs/tags/v2.0.0\n", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		rev, err := repo.ReadRevision(ctx, models.Omitted)
		require.NoError(t, err)
		assert.Equal(t, "latest-hash", rev.CommitHash)
		assert.Equal(t, "v2.0.0", rev.Version)
	})

	t.Run("read omitted (latest) without tag", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "ls-remote", "origin", "HEAD").
			Then("latest-hash\tHEAD\n", nil)
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "ls-remote", "origin").
			Then("some-other-hash\trefs/tags/v1.0.0\n", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		rev, err := repo.ReadRevision(ctx, models.Omitted)
		require.NoError(t, err)
		assert.Equal(t, "latest-hash", rev.CommitHash)
		assert.Equal(t, "latest-hash", rev.Version)
	})

	t.Run("read by commit hash", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		commitHash := "220e0db758f9ce96d9b1f457234616284530622b"

		// fetch call
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "fetch", "-f", "origin", "--depth=1", commitHash).
			Then("", nil)
		// get tags call
		consoleMock.RunCmdMock.When(ctx, "/tmp", "git", "ls-remote", "origin").
			Then("", nil)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		rev, err := repo.ReadRevision(ctx, models.RequestedVersion(commitHash))
		require.NoError(t, err)
		assert.Equal(t, commitHash, rev.CommitHash)
		assert.Equal(t, commitHash, rev.Version)
	})

	t.Run("version not found", func(t *testing.T) {
		consoleMock := mocks.NewConsoleMock(mc)
		consoleMock.RunCmdMock.Return("", assert.AnError)

		repo := &git.GitRepo{Console: consoleMock, CacheDir: "/tmp"}
		_, err := repo.ReadRevision(ctx, "v1.0.0")
		require.ErrorIs(t, err, models.ErrVersionNotFound)
	})
}
