package generate

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.redsock.ru/moti/internal/adapters/console"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
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

	return cmd
}

func (c Command) Action(cmd *cobra.Command, args []string) error {
	err := c.do(cmd, args)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate")
	}

	return nil
}

func (c Command) do(cmd *cobra.Command, _ []string) error {
	motiEnvironment, err := commands.GetEnvironment(cmd)
	if err != nil {
		return rerrors.Wrap(err)
	}

	dirWalker := fs.NewFSWalker(motiEnvironment.WorkDir, ".")

	app, err := buildCore(motiEnvironment, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	err = app.Generate(cmd.Context())
	if err != nil {
		return fmt.Errorf("generator.Generate: %w", err)
	}

	return nil
}

func buildCore(motiEnv commands.Env, dirWalker lockfile.DirWalker) (*Core, error) {
	lockFile, err := lockfile.New(dirWalker)
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening lockfile")
	}

	store := storage.New(motiEnv.MotiConfig.CachePath, lockFile)
	moduleCfg := moduleconfig.New()

	app := New(
		motiEnv.MotiConfig.Deps,
		&log.Logger,
		lo.Map(motiEnv.MotiConfig.Generate.Plugins, func(p config.Plugin, _ int) Plugin {
			return Plugin{
				Name:    p.Name,
				Out:     p.Out,
				Options: p.Opts,
			}
		}),
		Inputs{
			Dirs: lo.Filter(lo.Map(motiEnv.MotiConfig.Generate.Inputs, func(i config.Input, _ int) string {
				return i.Directory
			}), func(s string, _ int) bool {
				return s != ""
			}),
			InputGitRepos: lo.Filter(lo.Map(motiEnv.MotiConfig.Generate.Inputs, func(i config.Input, _ int) InputGitRepo {
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
		toolbox.Coalesce(motiEnv.MotiConfig.Generate.ProtoRoot, "."),
	)

	app.env = motiEnv

	return app, nil
}
