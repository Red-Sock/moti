package git

import (
	"context"
	"fmt"
	"path/filepath"

	"go.redsock.ru/moti/internal/core/models"
)

func (r *gitRepo) Archive(
	ctx context.Context, revision models.Revision, cacheDownloadPaths models.CacheDownloadPaths,
) error {
	absPath, err := filepath.Abs(cacheDownloadPaths.ArchiveFile)
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	params := []string{
		"archive", "--format=zip", revision.CommitHash, "-o", absPath, "*.proto",
	}

	if _, err := r.console.RunCmd(ctx, r.cacheDir, "git", params...); err != nil {
		return fmt.Errorf("utils.RunCmd: %w", err)
	}

	return nil
}
