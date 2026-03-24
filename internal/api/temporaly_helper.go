package api

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/samber/lo"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/adapters/console"
	"go.redsock.ru/moti/internal/adapters/go_git"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/core"
	"go.redsock.ru/moti/internal/rules"
)

func errExit(code int, msg string, args ...any) {
	log.Info().Msgf(msg, args...)
	os.Exit(code)
}

// deprecated: use internal/commands/generate.buildCore instead
func buildCore(_ context.Context, cfg config.Config, dirWalker core.DirWalker) (*core.Core, error) {
	lintRules, ignoreOnly, err := rules.New(cfg.Lint)
	if err != nil {
		return nil, rerrors.Wrap(err, "cfg.BuildLinterRules")
	}

	lockFile, err := lockfile.New(dirWalker)
	if err != nil {
		// TODO check no lock file
		return nil, rerrors.Wrap(err, "error opening lockfile")
	}

	store := storage.New(cfg.CachePath, lockFile)

	moduleCfg := moduleconfig.New()

	currentProjectGitWalker := go_git.New()

	breakingCheckConfig := core.BreakingCheckConfig{
		IgnoreDirs:    cfg.BreakingCheck.Ignore,
		AgainstGitRef: cfg.BreakingCheck.AgainstGitRef,
	}

	app := core.New(
		lintRules,
		cfg.Lint.Ignore,
		cfg.Deps,
		ignoreOnly,
		&log.Logger,
		lo.Map(cfg.Generate.Plugins, func(p config.Plugin, _ int) core.Plugin {
			return core.Plugin{
				Name:    p.Name,
				Out:     p.Out,
				Options: p.Opts,
			}
		}),
		core.Inputs{
			Dirs: lo.Filter(lo.Map(cfg.Generate.Inputs, func(i config.Input, _ int) string {
				return i.Directory
			}), func(s string, _ int) bool {
				return s != ""
			}),
			InputGitRepos: lo.Filter(lo.Map(cfg.Generate.Inputs, func(i config.Input, _ int) core.InputGitRepo {
				return core.InputGitRepo{
					URL:          i.GitRepo.URL,
					SubDirectory: i.GitRepo.SubDirectory,
					Out:          i.GitRepo.Out,
				}
			}), func(i core.InputGitRepo, _ int) bool {
				return i.URL != ""
			}),
		},
		console.New(),
		store,
		moduleCfg,
		lockFile,
		currentProjectGitWalker,
		breakingCheckConfig,
		cfg.Generate.ProtoRoot,
		cfg.Generate.GenerateOutDirs,
	)

	return app, nil
}
