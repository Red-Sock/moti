package generate

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/fs"
	"go.redsock.ru/moti/internal/models"
)

type Core struct {
	env commands.Env
}

func (c *Core) Generate(ctx context.Context) error {
	query := ProtocQuery{
		Imports: []string{
			toolbox.Coalesce(c.env.MotiConfig.Generate.ProtoRoot, "."),
		},
		Plugins: c.env.MotiConfig.Generate.Plugins,
	}

	err := mkdirForPluginsOut(c.env.MotiConfig.Generate.Plugins)
	if err != nil {
		return rerrors.Wrap(err, "mkdir for plugins failed")
	}

	for _, dep := range c.env.MotiConfig.Deps {
		modulePaths, err := c.getModulePath(dep)
		if err != nil {
			return rerrors.Wrap(err, "c.getModulePath")
		}

		query.Imports = append(query.Imports, modulePaths)
	}

	query, err = c.GenerateInputs(query)
	if err != nil {
		return rerrors.Wrap(err, "GenerateInputs")
	}

	command, args := query.Build()

	log.Info().
		Msg(command + " " + strings.Join(args, " \\\n           "))

	_, err = c.env.Console.RunCmd(ctx, c.env.WorkDir, command, args...)
	if err != nil {
		return rerrors.Wrap(err, "adapters.RunCmd")
	}

	return nil
}

func (c *Core) GenerateInputs(query ProtocQuery) (ProtocQuery, error) {
	for _, input := range c.env.MotiConfig.Generate.Inputs {
		moduleWalkerFunc := func(path string, err error) error {
			switch {
			case err != nil:
				return err
			case filepath.Ext(path) != ".proto":
				return nil
			}

			query.Files = append(query.Files, filepath.Join(input.Directory, path))

			return nil
		}

		if input.GitRepo.URL == "" {
			err := c.generateFromLocalFS(input, moduleWalkerFunc)
			if err != nil {
				return query, rerrors.Wrap(err, "generateFromLocalFS")
			}

			continue
		}

		err := c.generateFromGitRepo(&query, input, moduleWalkerFunc)
		if err != nil {
			return query, rerrors.Wrap(err, "generateFromGitRepo")
		}
	}

	return query, nil
}

func (c *Core) generateFromLocalFS(input config.Input, walker func(path string, err error) error) error {
	fsWalker := fs.NewFSWalker(input.Directory, "")

	err := fsWalker.WalkDir(walker)
	if err != nil {
		return rerrors.Wrap(err, "fsWalker.WalkDir")
	}

	return nil
}

func (c *Core) generateFromGitRepo(query *ProtocQuery, input config.Input,
	walker func(path string, err error) error) error {
	module := models.NewModule(input.GitRepo.URL)

	isInstalled, err := c.env.Storage.IsModuleInstalled(module)
	if err != nil {
		return rerrors.Wrap(err, "c.isModuleInstalled")
	}

	if !isInstalled {
		return rerrors.Wrap(models.ErrModuleNotInstalled, module)
	}

	modulePaths, err := c.getModulePath(module.Name)
	if err != nil {
		return rerrors.Wrap(err, "c.getModulePath")
	}

	fsWalker := fs.NewFSWalker(modulePaths, input.GitRepo.SubDirectory)

	gitGenerateCb := func(path string, err error) error {
		err = walker(path, err)
		if err != nil {
			return err
		}

		query.Imports = append(query.Imports, modulePaths)

		return nil
	}

	err = fsWalker.WalkDir(gitGenerateCb)
	if err != nil {
		return rerrors.Wrap(err, "fsWalker.WalkDir")
	}

	return nil
}

func (c *Core) getModulePath(requestedDependency string) (string, error) {
	module := models.NewModule(requestedDependency)

	isInstalled, err := c.env.Storage.IsModuleInstalled(module)
	if err != nil {
		return "", rerrors.Wrap(err, "h.storage.IsModuleInstalled")
	}

	if !isInstalled {
		return "", rerrors.Wrap(models.ErrModuleNotInstalled, module.Name)
	}

	lockFileInfo, err := c.env.LockFile.Read(module.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "lockFile.Read")
	}

	return c.env.Storage.GetInstallDir(module.Name, lockFileInfo.Version), nil
}

func toUniqueMap(imports []string) map[string]struct{} {
	uniqueImports := make(map[string]struct{}, len(imports))

	for _, imp := range imports {
		if imp == "" {
			continue
		}

		uniqueImports[imp] = struct{}{}
	}

	return uniqueImports
}

func mkdirForPluginsOut(plugins []config.Plugin) error {
	for _, plug := range plugins {
		if filepath.IsAbs(plug.Out) {
			continue
		}

		err := os.MkdirAll(plug.Out, storage.DirPerm)
		if err != nil {
			return rerrors.Wrap(err, "os.MkdirAll")
		}
	}

	return nil
}
