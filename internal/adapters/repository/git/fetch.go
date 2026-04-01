package git

import (
	"context"
	"fmt"

	"go.redsock.ru/moti/internal/models"
)

func (r *GitRepo) Fetch(ctx context.Context, revision models.Revision) error {
	_, err := r.Console.RunCmd(
		ctx, r.CacheDir,
		"git", "fetch",
		"-f", "origin",
		"--depth=1",
		revision.CommitHash,
	)
	if err != nil {
		return fmt.Errorf("adapters.RunCmd (fetch): %w", err)
	}

	return nil
}
