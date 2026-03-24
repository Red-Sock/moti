package moduleconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"

	"gopkg.in/yaml.v3"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/core/models"
)

type bufWork struct {
	Directories []string `yaml:"directories"`
}

const (
	bufWorkFile = "buf.work.yaml"
)

func readBufWork(ctx context.Context, repo repository.Repo, revision models.Revision) (bufWork, error) {
	content, err := repo.ReadFile(ctx, revision, bufWorkFile)
	if err != nil {
		if errors.Is(err, models.ErrFileNotFound) {
			log.Debug().Msg("buf config not found")
			return bufWork{}, nil
		}
		return bufWork{}, fmt.Errorf("repo.ReadFile: %w", err)
	}

	buf := bufWork{}
	if err := yaml.NewDecoder(strings.NewReader(content)).Decode(&buf); err != nil {
		return bufWork{}, fmt.Errorf("yaml.NewDecoder: %w", err)
	}

	return buf, nil
}
