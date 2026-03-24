package moduleconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"

	"gopkg.in/yaml.v3"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/config/default_consts"
	"go.redsock.ru/moti/internal/core/models"
)

// Config is the configuration of moti.
// FIXME: do not duplicate of struct
// but if now will import from config -> cycles deps
type motiConfig struct {
	// Deps is the dependencies repositories
	Deps []string `json:"deps" yaml:"deps"`
}

// readmoti read moti's config from repository
func readmoti(ctx context.Context, repo repository.Repo, revision models.Revision) ([]models.Module, error) {
	content, err := repo.ReadFile(ctx, revision, default_consts.DefaultConfigFileName)
	if err != nil {
		if errors.Is(err, models.ErrFileNotFound) {
			log.Debug().Msg("moti config not found")
			return nil, nil
		}
		return nil, fmt.Errorf("repo.ReadFile: %w", err)
	}

	moti := &motiConfig{}
	if err := yaml.NewDecoder(strings.NewReader(content)).Decode(&moti); err != nil {
		return nil, fmt.Errorf("yaml.NewDecoder: %w", err)
	}

	modules := make([]models.Module, 0, len(moti.Deps))
	for _, dep := range moti.Deps {
		module := models.NewModule(dep)
		modules = append(modules, module)
	}

	return modules, nil
}
