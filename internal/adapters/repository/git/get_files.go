package git

import (
	"context"
	"fmt"
	"strings"

	"go.redsock.ru/moti/internal/models"
)

func (r *gitRepo) GetFiles(ctx context.Context, revision models.Revision, dirs ...string) ([]string, error) {
	params := make([]string, 0, 3+len(dirs))

	params = append(params, "ls-tree", "-r", revision.CommitHash)
	params = append(params, dirs...)

	res, err := r.console.RunCmd(ctx, r.cacheDir, "git", params...)
	if err != nil {
		return nil, fmt.Errorf("utils.RunCmd: %w", err)
	}

	stats := strings.Split(res, "\n")

	files := make([]string, 0, len(stats))
	for _, stat := range stats {
		stat := stat

		statFields := strings.Fields(stat)
		if len(statFields) != 4 {
			continue
		}

		files = append(files, statFields[3])
	}

	return files, nil
}
