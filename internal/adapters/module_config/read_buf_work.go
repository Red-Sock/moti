package moduleconfig

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"gopkg.in/yaml.v3"

	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/models"
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
			log.Debug().
				Msg("buf config not found")

			return bufWork{}, nil
		}

		return bufWork{}, rerrors.Wrap(err, "repo.ReadFile")
	}

	buf := bufWork{}

	err = yaml.NewDecoder(strings.NewReader(content)).Decode(&buf)
	if err != nil {
		return bufWork{}, rerrors.Wrap(err, "yaml.NewDecoder")
	}

	return buf, nil
}
