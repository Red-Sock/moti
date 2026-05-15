package config

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CachePath string `json:"cache_path" yaml:"cache_path"`

	Deps []string `json:"deps" yaml:"deps"`

	Replace []Replace `json:"replace" yaml:"replace"`

	Binaries Binaries `json:"binaries" yaml:"binaries"`

	Generate []Generate `json:"generate" yaml:"generate"`
}

func Read(cfgPath string) (Config, error) {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrFileNotFound
		}

		return Config{}, rerrors.Wrap(err, "error opening config file")
	}

	defer func() {
		err = cfgFile.Close()
		if err != nil {
			log.Debug().Err(err).Str("filepath", cfgPath).Msg("error closing config file")
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

	cfg.CachePath = toolbox.Coalesce(cfg.CachePath, "proto_modules")

	cfgDir := filepath.Dir(cfgPath)
	for i, r := range cfg.Replace {
		if !filepath.IsAbs(r.New) {
			cfg.Replace[i].New = filepath.Join(cfgDir, r.New)
		}
	}

	return cfg, nil
}

func ReadOrDie(cmd *cobra.Command) Config {
	configPath, _ := cmd.Flags().GetString(ConfigFlag)

	cfg, err := Read(configPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("filepath", configPath).
			Msg("error reading config file")
	}

	return cfg
}

func (c Config) BuildPATH(workDir string) string {
	if c.Binaries.BinDir == "" {
		return ""
	}

	return PATHPrefix + path.Join(workDir, c.Binaries.BinDir) + ":$PATH"
}

func (c Config) BuildGOBIN(workDir string) string {
	if c.Binaries.BinDir == "" {
		return ""
	}

	return GOBINPrefix + path.Join(workDir, c.Binaries.BinDir)
}
