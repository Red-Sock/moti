package moduleconfig

import (
	"context"
	"fmt"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/core/models"
)

// ReadFromRepo  Read and return module's config from repository
func (c *ModuleConfig) ReadFromRepo(
	ctx context.Context, repo repository.Repo, revision models.Revision,
) (models.ModuleConfig, error) {
	buf, err := readBufWork(ctx, repo, revision)
	if err != nil {
		return models.ModuleConfig{}, fmt.Errorf("readBufWork: %w", err)
	}

	modules, err := readmoti(ctx, repo, revision)
	if err != nil {
		return models.ModuleConfig{}, fmt.Errorf("readmoti: %w", err)
	}

	moduleConfig := models.ModuleConfig{
		Directories:  buf.Directories,
		Dependencies: modules,
	}
	return moduleConfig, nil
}
