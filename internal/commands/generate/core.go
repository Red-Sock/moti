package generate

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/adapters/console"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/core/models"
	"go.redsock.ru/moti/internal/fs/fs"
)

type (
	Plugin struct {
		Name    string
		Out     string
		Options map[string]string
	}

	InputGitRepo struct {
		URL          string
		SubDirectory string
		Out          string
	}

	Inputs struct {
		Dirs          []string
		InputGitRepos []InputGitRepo
	}

	Query struct {
		Compiler string
		Files    []string
		Imports  []string
		Plugins  []Plugin
	}
)

func (q Query) Build() (command string, args []string) {
	command = q.Compiler

	for _, imp := range slices.Sorted(maps.Keys(q.getUniqueImports())) {
		args = append(args, "-I "+imp)
	}

	for _, plug := range q.Plugins {
		arg := "--" + plug.Name + "_out="

		var opts []string
		for k, v := range plug.Options {
			if v != "" {
				opts = append(opts, fmt.Sprintf("%s=%s", k, v))
			} else {
				opts = append(opts, k)
			}
		}

		if len(opts) > 0 {
			arg += strings.Join(opts, ",") + ":"
		}

		arg += plug.Out
		args = append(args, arg)
	}

	uniqueProtoFileDirs := make(map[string]struct{})
	for _, file := range q.Files {
		uniqueProtoFileDirs[filepath.Dir(file)] = struct{}{}
	}

	for file := range uniqueProtoFileDirs {
		args = append(args, file+"/*.proto")
	}

	return command, args
}

func (q Query) buildWithResponseFile(command string, args []string, _ string) (string, []string, string, error) {
	tmpFile, err := os.CreateTemp("", "protoc-args-*.rsp")
	if err != nil {
		return command, args, "", err
	}
	tmpFileName := tmpFile.Name()

	for _, arg := range args {
		if strings.ContainsAny(arg, " \t\n\r\"'") {
			arg = "\"" + strings.ReplaceAll(arg, "\"", "\\\"") + "\""
		}
		_, err = tmpFile.WriteString(arg + "\n")
		if err != nil {
			_ = tmpFile.Close()
			return command, args, tmpFileName, err
		}
	}

	err = tmpFile.Close()
	if err != nil {
		return command, args, tmpFileName, err
	}

	return command, []string{"@" + tmpFileName}, tmpFileName, nil
}

func (q Query) getUniqueImports() map[string]struct{} {
	uniqueImports := make(map[string]struct{}, len(q.Imports)+len(q.Files))

	for _, imp := range q.Imports {
		if imp == "" {
			continue
		}
		uniqueImports[imp] = struct{}{}
	}

	//for _, file := range q.Files {
	//	uniqueImports[filepath.Dir(file)] = struct{}{}
	//}

	return uniqueImports
}

type Core struct {
	env commands.Env

	deps         []string
	logger       *zerolog.Logger
	plugins      []Plugin
	inputs       Inputs
	console      console.Console
	storage      storage.IStorage
	moduleConfig moduleconfig.IModuleConfig
	lockFile     lockfile.ILockFile
	protoRoot    string
}

func New(
	deps []string,
	logger *zerolog.Logger,
	plugins []Plugin,
	inputs Inputs,
	console console.Console,
	storage storage.IStorage,
	moduleConfig moduleconfig.IModuleConfig,
	lockFile lockfile.ILockFile,
	protoRoot string,
) *Core {
	return &Core{
		deps:         deps,
		logger:       logger,
		plugins:      plugins,
		inputs:       inputs,
		console:      console,
		storage:      storage,
		moduleConfig: moduleConfig,
		lockFile:     lockFile,
		protoRoot:    protoRoot,
	}
}

func (c *Core) Generate(ctx context.Context) error {
	q := Query{
		Compiler: "protoc",
		Imports: []string{
			c.protoRoot,
		},
		Plugins: c.plugins,
	}

	for _, dep := range c.deps {
		modulePaths, err := c.getModulePath(dep)
		if err != nil {
			return rerrors.Wrap(err, "c.getModulePath")
		}

		q.Imports = append(q.Imports, modulePaths)
	}

	for _, plug := range q.Plugins {
		if filepath.IsAbs(plug.Out) {
			continue
		}

		err := os.MkdirAll(plug.Out, os.ModePerm)
		if err != nil {
			return rerrors.Wrap(err, "os.MkdirAll")
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
			modulePaths, err := c.getModulePath(module.Name)
			if err != nil {
				return rerrors.Wrap(err, "c.getModulePath")
			}

			fsWalker := fs.NewFSWalker(modulePaths, repo.SubDirectory)

			err = fsWalker.WalkDir(gitGenerateCb(modulePaths))
			if err != nil {
				return rerrors.Wrap(err, "fsWalker.WalkDir1: %w")
			}

			continue
		}

		return fmt.Errorf("module %s is not installed. Please run `moti install` first", module.Name)
	}

	for _, inp := range c.inputs.Dirs {
		fsWalker := fs.NewFSWalker(inp, "")
		err := fsWalker.WalkDir(func(path string, err error) error {
			switch {
			case err != nil:
				return err
			case ctx.Err() != nil:
				return ctx.Err()
			case filepath.Ext(path) != ".proto":
				return nil
			}

			q.Files = append(q.Files, filepath.Join(inp, path))
			//q.Imports = append(q.Imports, inp)

			return nil
		})
		if err != nil {
			return rerrors.Wrap(err, "fsWalker.WalkDir")
		}
	}

	command, args := q.Build()

	_, err := c.console.RunCmd(ctx, c.env.WorkDir, command, args...)
	if err != nil {
		return rerrors.Wrap(err, "adapters.RunCmd")
	}

	return nil
}

func (c *Core) getModulePath(requestedDependency string) (string, error) {
	module := models.NewModule(requestedDependency)

	isInstalled, err := c.storage.IsModuleInstalled(module)
	if err != nil {
		return "", rerrors.Wrap(err, "h.storage.IsModuleInstalled")
	}

	if !isInstalled {
		return "", rerrors.Wrap(err, "module %s is not installed. Please run `moti install` first", module.Name)
	}

	lockFileInfo, err := c.lockFile.Read(module.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "lockFile.Read")
	}

	installedPath := c.storage.GetInstallDir(module.Name, lockFileInfo.Version)

	return installedPath, nil
}
