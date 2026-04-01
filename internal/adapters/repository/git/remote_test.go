package git_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/mocks"
)

func TestGitRepo_New(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	t.Run("new repository initialization", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "moti-git-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		consoleMock := mocks.NewConsoleMock(mc)
		// git init
		consoleMock.RunCmdMock.When(ctx, tempDir, "git", "init", "--bare").Then("", nil)
		// git remote add
		// Note: GetRemote will be called, we need to mock it if it's not a real URL.
		// Since we're using a real URL in the test below, it might try to reach it.
		// Let's use a httptest server.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `<html><head><meta name="go-import" content="example.com/repo git https://github.com/example/repo"></head></html>`)
		}))
		defer server.Close()

		remote := server.URL // keep schema

		consoleMock.RunCmdMock.When(ctx, tempDir, "git", "remote", "add", "origin", "https://github.com/example/repo").Then("", nil)

		repo, err := git.New(ctx, remote, tempDir, consoleMock)
		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("existing repository", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "moti-git-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		err = os.MkdirAll(filepath.Join(tempDir, "objects"), 0755)
		require.NoError(t, err)

		consoleMock := mocks.NewConsoleMock(mc)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `<html></html>`)
		}))
		defer server.Close()
		remote := server.URL

		repo, err := git.New(ctx, remote, tempDir, consoleMock)
		require.NoError(t, err)
		assert.NotNil(t, repo)
	})
}
