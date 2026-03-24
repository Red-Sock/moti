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

// Config is the configuration of protopack.
// FIXME: do not duplicate of struct
// but if now will import from config -> cycles deps
type protopackConfig struct {
	// Deps is the dependencies repositories
	Deps []string `json:"deps" yaml:"deps"`
}

// readprotopack read protopack's config from repository
func readProtoPack(ctx context.Context, repo repository.Repo, revision models.Revision) ([]models.Module, error) {
	content, err := repo.ReadFile(ctx, revision, default_consts.DefaultConfigFileName)
	if err != nil {
		if errors.Is(err, models.ErrFileNotFound) {
			log.Debug().Msg("protopack config not found")
			return nil, nil
		}
		return nil, fmt.Errorf("repo.ReadFile: %w", err)
	}

	protopack := &protopackConfig{}
	if err := yaml.NewDecoder(strings.NewReader(content)).Decode(&protopack); err != nil {
		return nil, fmt.Errorf("yaml.NewDecoder: %w", err)
	}

	modules := make([]models.Module, 0, len(protopack.Deps))
	for _, dep := range protopack.Deps {
		module := models.NewModule(dep)
		modules = append(modules, module)
	}

	return modules, nil
}
