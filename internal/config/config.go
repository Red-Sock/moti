package config

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CachePath string `json:"cache_path" yaml:"cache_path"`

	Deps []string `json:"deps" yaml:"deps"`

	Generate Generate `json:"generate" yaml:"generate"`
}

func Read(filepath string) (Config, error) {
	cfgFile, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrFileNotFound
		}

		return Config{}, rerrors.Wrap(err, "error opening config file")
	}

	defer func() {
		err = cfgFile.Close()
		if err != nil {
			log.Debug().Err(err).Str("filepath", filepath).Msg("error closing config file")
		}
	}()

	buf, err := io.ReadAll(cfgFile)
	if err != nil {
		return Config{}, rerrors.Wrap(err, "error reading config file")
	}

	cfg := Config{}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return Config{}, rerrors.Wrap(err, "error parsing config file from yaml")
	}

	return cfg, nil
}
