package commands

import (
	"os"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/flags"
)

type Env struct {
	WorkDir    string
	MotiConfig config.Config
}

func GetEnvironment(cmd *cobra.Command) (Env, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return Env{}, rerrors.Wrap(err, "os.Getwd")
	}

	configPath, _ := cmd.Flags().GetString(flags.Config)
	motiCfg, err := config.Read(configPath)
	if err != nil {
		return Env{}, rerrors.Wrap(err, "config.Read")
	}

	return Env{
		WorkDir:    workingDir,
		MotiConfig: motiCfg,
	}, nil
}
