package git

import (
	"context"
	"path/filepath"

	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/models"
)

func (r *GitRepo) Archive(
	ctx context.Context, revision models.Revision, cacheDownloadPaths models.CacheDownloadPaths) error {
	absPath, err := filepath.Abs(cacheDownloadPaths.ArchiveFile)
	if err != nil {
		return rerrors.Wrap(err, "filepath.Abs")
	}

	params := []string{
		"archive",
		"--format=zip",
		revision.CommitHash,
		"-o", absPath,
		"*.proto",
	}

	_, err = r.Console.RunCmd(ctx, r.CacheDir, "git", params...)
	if err != nil {
		return rerrors.Wrap(err, "utils.RunCmd")
	}

	return nil
}
