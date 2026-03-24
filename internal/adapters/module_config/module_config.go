package moduleconfig

import (
	"context"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/core/models"
)

type (
	// ModuleConfig implement module config logic such as buf dirs config etc
	ModuleConfig struct {
	}

	IModuleConfig interface {
		ReadFromRepo(ctx context.Context, repo repository.Repo, revision models.Revision) (models.ModuleConfig, error)
	}
)

func New() *ModuleConfig {
	return &ModuleConfig{}
}
