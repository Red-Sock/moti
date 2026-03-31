package moduleconfig

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/models"
)

type motiConfig struct {
	Deps []string `json:"deps" yaml:"deps"`
}

func readMoti(ctx context.Context, repo repository.Repo, revision models.Revision) ([]models.Module, error) {
	content, err := repo.ReadFile(ctx, revision, config.DefaultConfigFilePath)
	if err != nil {
		if errors.Is(err, models.ErrFileNotFound) {
			log.Debug().
				Msg("moti config not found")

			return nil, nil
		}

		return nil, rerrors.Wrap(err, "repo.ReadFile")
	}

	moti := &motiConfig{}

	err = yaml.NewDecoder(strings.NewReader(content)).Decode(&moti)
	if err != nil {
		return nil, rerrors.Wrap(err, "yaml.NewDecoder")
	}

	modules := make([]models.Module, 0, len(moti.Deps))
	for _, dep := range moti.Deps {
		module := models.NewModule(dep)
		modules = append(modules, module)
	}

	return modules, nil
}
