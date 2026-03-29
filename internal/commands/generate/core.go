package generate

import (
	"context"
	"os"
	"path"
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
	Env    commands.Env
	Walker IWalker
}

//go:generate minimock -i IWalker -o ../../mocks -g -s "_mock.go"
type IWalker interface {
	WalkDir(root, path string, callback func(path string, err error) error) error
}

type fsWalker struct{}

func (f *fsWalker) WalkDir(root, path string, callback func(path string, err error) error) error {
	w := fs.NewFSWalker(root, path)
	return w.WalkDir(callback)
}

func (c *Core) Generate(ctx context.Context) error {
	if c.Walker == nil {
		c.Walker = &fsWalker{}
	}

	for _, genCfg := range c.Env.MotiConfig.Generate {
		err := mkdirForPluginsOut(genCfg.Plugins)
		if err != nil {
			return rerrors.Wrap(err, "mkdir for plugins failed")
		}

		for _, input := range genCfg.Inputs {
			query := ProtocQuery{
				Imports: []string{
					c.getFirstInput(input),
				},
				Plugins: genCfg.Plugins,
			}

			for _, dep := range c.Env.MotiConfig.Deps {
				modulePaths, err := c.getModulePath(dep)
				if err != nil {
					return rerrors.Wrap(err, "c.getModulePath")
				}

				query.Imports = append(query.Imports, modulePaths)
			}

			if input.GitRepo.URL == "" {
				err = c.generateFromLocalFS(&query, input)
				if err != nil {
					return rerrors.Wrap(err, "generateFromLocalFS")
				}
			} else {
				err = c.generateFromGitRepo(&query, input)
				if err != nil {
					return rerrors.Wrap(err, "generateFromGitRepo")
				}
			}

			command, args := query.Build()

			log.Info().
				Msg(command + " " + strings.Join(args, " \\\n           "))

			_, err = c.Env.Console.RunCmd(ctx, c.Env.WorkDir, command, args...)
			if err != nil {
				return rerrors.Wrap(err, "adapters.RunCmd")
			}
		}
	}

	return nil
}

func (c *Core) generateFromLocalFS(query *ProtocQuery, input config.Input) error {
	walker := func(path string, err error) error {
		isProto, err := isContainingProto(path, err)
		if err != nil {
			return rerrors.Wrap(err, "isContainingProto")
		}

		if isProto {
			query.Files = append(query.Files, filepath.Join(input.Directory, path))
		}

		return nil
	}

	err := c.Walker.WalkDir(input.Directory, "", walker)
	if err != nil {
		return rerrors.Wrap(err, "Walker.WalkDir")
	}

	return nil
}

func (c *Core) generateFromGitRepo(query *ProtocQuery, input config.Input) error {
	module := models.NewModule(input.GitRepo.URL)

	isInstalled, err := c.Env.Storage.IsModuleInstalled(module)
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

	gitGenerateCb := func(path string, err error) error {
		containsProto, err := isContainingProto(path, err)
		if err != nil {
			return err
		}

		if containsProto {
			moduleProtoPath := filepath.Join(modulePaths, path)

			query.Files = append(query.Files, moduleProtoPath)
			query.Imports = append(query.Imports, filepath.Dir(moduleProtoPath))
		}

		return nil
	}

	err = c.Walker.WalkDir(modulePaths, input.GitRepo.SubDirectory, gitGenerateCb)
	if err != nil {
		return rerrors.Wrap(err, "Walker.WalkDir")
	}

	return nil
}

func (c *Core) getModulePath(requestedDependency string) (string, error) {
	module := models.NewModule(requestedDependency)

	isInstalled, err := c.Env.Storage.IsModuleInstalled(module)
	if err != nil {
		return "", rerrors.Wrap(err, "h.storage.IsModuleInstalled")
	}

	if !isInstalled {
		return "", rerrors.Wrap(models.ErrModuleNotInstalled, module.Name)
	}

	lockFileInfo, err := c.Env.LockFile.Read(module.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "lockFile.Read")
	}

	return c.Env.Storage.GetInstallDir(module.Name, lockFileInfo.Version), nil
}

func (c *Core) getFirstInput(inp config.Input) string {
	if inp.GitRepo.URL == "" {
		return toolbox.Coalesce(inp.Directory, ".")
	}

	modulePath, err := c.getModulePath(inp.GitRepo.URL)
	if err != nil {
		log.Panic().
			Err(err).
			Str("url", inp.GitRepo.URL).
			Msg("Error getting module path for git repo")
	}

	return path.Join(modulePath, inp.GitRepo.SubDirectory)
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

func isContainingProto(path string, err error) (bool, error) {
	switch {
	case err != nil:
		return false, err
	case filepath.Ext(path) != ".proto":
		return false, nil
	}

	return true, nil
}
