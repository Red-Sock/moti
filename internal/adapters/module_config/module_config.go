package moduleconfig

import (
	"context"

	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/models"
)

type ModuleConfig struct {
}

type IModuleConfig interface {
	ReadFromRepo(ctx context.Context, repo repository.Repo, revision models.Revision) (models.ModuleConfig, error)
}

func New() *ModuleConfig {
	return &ModuleConfig{}
}

func (c *ModuleConfig) ReadFromRepo(
	ctx context.Context, repo repository.Repo, rev models.Revision) (models.ModuleConfig, error) {
	buf, err := readBufWork(ctx, repo, rev)
	if err != nil {
		return models.ModuleConfig{}, rerrors.Wrap(err, "readBufWork")
	}

	modules, err := readMoti(ctx, repo, rev)
	if err != nil {
		return models.ModuleConfig{}, rerrors.Wrap(err, "read moti")
	}

	moduleConfig := models.ModuleConfig{
		Directories:  buf.Directories,
		Dependencies: modules,
	}

	return moduleConfig, nil
}
