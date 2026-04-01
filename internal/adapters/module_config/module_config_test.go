package moduleconfig

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.redsock.ru/moti/internal/models"
)

type mockRepo struct {
	files map[string]string
}

func (m *mockRepo) ReadFile(ctx context.Context, revision models.Revision, fileName string) (string, error) {
	content, ok := m.files[fileName]
	if !ok {
		return "", models.ErrFileNotFound
	}
	return content, nil
}

func (m *mockRepo) Archive(ctx context.Context, revision models.Revision, cacheDownloadPaths models.CacheDownloadPaths) error {
	return nil
}

func (m *mockRepo) ReadRevision(ctx context.Context, requestedVersion models.RequestedVersion) (models.Revision, error) {
	return models.Revision{}, nil
}

func (m *mockRepo) Fetch(ctx context.Context, revision models.Revision) error {
	return nil
}

func TestModuleConfig_ReadFromRepo(t *testing.T) {
	ctx := context.Background()
	rev := models.Revision{CommitHash: "hash"}
	mc := New()

	t.Run("success with both configs", func(t *testing.T) {
		repo := &mockRepo{
			files: map[string]string{
				"buf.work.yaml": "directories:\n  - proto\n  - api\n",
				"moti.yaml":     "deps:\n  - github.com/user/repo@v1.0.0\n",
			},
		}

		res, err := mc.ReadFromRepo(ctx, repo, rev)
		require.NoError(t, err)

		assert.Equal(t, []string{"proto", "api"}, res.Directories)
		require.Len(t, res.Dependencies, 1)
		assert.Equal(t, "github.com/user/repo", res.Dependencies[0].Name)
		assert.Equal(t, models.RequestedVersion("v1.0.0"), res.Dependencies[0].Version)
	})

	t.Run("no configs", func(t *testing.T) {
		repo := &mockRepo{files: map[string]string{}}

		res, err := mc.ReadFromRepo(ctx, repo, rev)
		require.NoError(t, err)

		assert.Empty(t, res.Directories)
		assert.Empty(t, res.Dependencies)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		repo := &mockRepo{
			files: map[string]string{
				"moti.yaml": "invalid: yaml: :",
			},
		}

		_, err := mc.ReadFromRepo(ctx, repo, rev)
		assert.Error(t, err)
	})
}
