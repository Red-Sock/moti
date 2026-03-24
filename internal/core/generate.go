package core

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.redsock.ru/moti/internal/core/models"
	"go.redsock.ru/moti/internal/fs/fs"
)

const (
	protocCompiler = "protoc"

	defaultCompiler = protocCompiler
)

// Generate generates files.
// deprecated: use internal/commands/generate.Core instead
func (c *Core) Generate(ctx context.Context, root, directory string) error {
	q := Query{
		Compiler: defaultCompiler,
		Imports: []string{
			toolbox.Coalesce(c.protoRoot, root),
		},
		Plugins: c.plugins,
	}

	for _, dep := range c.deps {
		modulePaths, err := c.getModulePath(ctx, dep)
		if err != nil {
			return rerrors.Wrap(err, "g.moduleReflect.GetModulePath")
		}

		q.Imports = append(q.Imports, modulePaths)
	}

	if c.generateOutDirs {
		for _, plug := range q.Plugins {
			if filepath.IsAbs(plug.Out) {
				continue
			}

			err := os.MkdirAll(plug.Out, os.ModePerm)
			if err != nil {
				return rerrors.Wrap(err, "os.MkdirAll")
			}
		}
	}

	for _, repo := range c.inputs.InputGitRepos {
		module := models.NewModule(repo.URL)

		isInstalled, err := c.storage.IsModuleInstalled(module)
		if err != nil {
			return rerrors.Wrap(err, "c.isModuleInstalled")
		}

		gitGenerateCb := func(modulePaths string) func(path string, err error) error {
			return func(path string, err error) error {
				switch {
				case err != nil:
					return err
				case ctx.Err() != nil:
					return ctx.Err()
				case filepath.Ext(path) != ".proto":
					return nil
				}

				q.Files = append(q.Files, path)
				q.Imports = append(q.Imports, modulePaths)

				return nil
			}
		}

		if isInstalled {
			modulePaths, err := c.getModulePath(ctx, module.Name)
			if err != nil {
				return rerrors.Wrap(err, "g.moduleReflect.GetModulePath")
			}

			fsWalker := fs.NewFSWalker(modulePaths, repo.SubDirectory)

			err = fsWalker.WalkDir(gitGenerateCb(modulePaths))
			if err != nil {
				return rerrors.Wrap(err, "fsWalker.WalkDir1: %w")
			}

			continue
		}

		err = c.InstallPackage(ctx, module)
		if err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().Str("dependency", module.Name).Str("version", string(module.Version)).Msg("Version not found")

				return rerrors.Wrap(err, "models.ErrVersionNotFound")
			}

			return rerrors.Wrap(err, "c.Get")
		}

		modulePaths, err := c.getModulePath(ctx, module.Name)
		if err != nil {
			return rerrors.Wrap(err, "g.moduleReflect.GetModulePath")
		}

		fsWalker := fs.NewFSWalker(modulePaths, repo.SubDirectory)
		err = fsWalker.WalkDir(gitGenerateCb(modulePaths))
		if err != nil {
			return rerrors.Wrap(err, "fsWalker.WalkDir")
		}
	}

	fsWalker := fs.NewFSWalker(directory, "")
	err := fsWalker.WalkDir(func(path string, err error) error {
		switch {
		case err != nil:
			return err
		case ctx.Err() != nil:
			return ctx.Err()
		case filepath.Ext(path) != ".proto":
			return nil
		case shouldIgnore(path, c.inputs.Dirs):
			c.logger.Debug().Str("path", path).Msg("ignore")

			return nil
		}

		q.Files = append(q.Files, path)

		return nil
	})
	if err != nil {
		return rerrors.Wrap(err, "fsWalker.WalkDir")
	}

	cmd := q.build()

	log.Info().Msg("Run command")
	println(cmd)

	_, err = c.console.RunCmd(ctx, root, cmd)
	if err != nil {
		return rerrors.Wrap(err, "adapters.RunCmd")
	}

	return nil
}

func shouldIgnore(path string, dirs []string) bool {
	if len(dirs) == 0 {
		return false
	}

	for _, dir := range dirs {
		if strings.HasPrefix(path, dir) {
			return false
		}
	}

	return true
}

func (c *Core) getModulePath(ctx context.Context, requestedDependency string) (string, error) {
	module := models.NewModule(requestedDependency)

	isInstalled, err := c.storage.IsModuleInstalled(module)
	if err != nil {
		return "", rerrors.Wrap(err, "h.storage.IsModuleInstalled")
	}

	if !isInstalled {
		err = c.InstallPackage(ctx, module)
		if err != nil {
			return "", rerrors.Wrap(err, "h.mod.Get")
		}
	}

	lockFileInfo, err := c.lockFile.Read(module.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "lockFile.Read")
	}

	installedPath := c.storage.GetInstallDir(module.Name, lockFileInfo.Version)

	return installedPath, nil
}
