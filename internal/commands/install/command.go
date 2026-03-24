package install

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/adapters/console"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/fs/fs"
)

type Command struct{}

func (c Command) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "install dependencies",
		Long:    "install dependencies specified in moti.yaml",
		RunE:    c.Action,
	}

	return cmd
}

func (c Command) Action(cmd *cobra.Command, args []string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	dirWalker := fs.NewFSWalker(workingDir, ".")

	app, err := buildCore(*cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	err = app.Install(cmd.Context(), cfg.Deps)
	if err != nil {
		return fmt.Errorf("install: %w", err)
	}

	return nil
}

func buildCore(cfg config.Config, dirWalker lockfile.DirWalker) (*Core, error) {
	lockFile, err := lockfile.New(dirWalker)
	if err != nil {
		return nil, fmt.Errorf("error opening lockfile: %w", err)
	}

	store := storage.New(cfg.CachePath, lockFile)
	moduleCfg := moduleconfig.New()

	app := New(
		&log.Logger,
		console.New(),
		store,
		moduleCfg,
		lockFile,
	)

	return app, nil
}
