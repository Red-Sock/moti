package generate

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

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
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "generate code from proto files",
		Long:    "generate code from proto files",
		RunE:    c.Action,
	}

	cmd.Flags().StringP("path", "p", ".", "set path to directory with proto files")
	_ = cmd.MarkFlagRequired("path")

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

	dir, _ := cmd.Flags().GetString("path")
	err = app.Generate(cmd.Context(), ".", dir)
	if err != nil {
		return fmt.Errorf("generator.Generate: %w", err)
	}

	return nil
}

func buildCore(cfg config.Config, dirWalker lockfile.DirWalker) (*Core, error) {
	lockFile, err := lockfile.New(dirWalker)
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening lockfile")
	}

	store := storage.New(cfg.CachePath, lockFile)
	moduleCfg := moduleconfig.New()

	app := New(
		cfg.Deps,
		&log.Logger,
		lo.Map(cfg.Generate.Plugins, func(p config.Plugin, _ int) Plugin {
			return Plugin{
				Name:    p.Name,
				Out:     p.Out,
				Options: p.Opts,
			}
		}),
		Inputs{
			Dirs: lo.Filter(lo.Map(cfg.Generate.Inputs, func(i config.Input, _ int) string {
				return i.Directory
			}), func(s string, _ int) bool {
				return s != ""
			}),
			InputGitRepos: lo.Filter(lo.Map(cfg.Generate.Inputs, func(i config.Input, _ int) InputGitRepo {
				return InputGitRepo{
					URL:          i.GitRepo.URL,
					SubDirectory: i.GitRepo.SubDirectory,
					Out:          i.GitRepo.Out,
				}
			}), func(i InputGitRepo, _ int) bool {
				return i.URL != ""
			}),
		},
		console.New(),
		store,
		moduleCfg,
		lockFile,
		cfg.Generate.ProtoRoot,
		cfg.Generate.GenerateOutDirs,
	)

	return app, nil
}
