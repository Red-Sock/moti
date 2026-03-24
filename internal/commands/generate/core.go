package generate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.redsock.ru/moti/internal/adapters/console"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
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

func (q Query) build() string {
	var buf bytes.Buffer
	buf.WriteString(q.Compiler)

	for _, imp := range q.Imports {
		buf.WriteString(" -I")
		buf.WriteString(imp)
	}

	for _, plug := range q.Plugins {
		buf.WriteString(" --")
		buf.WriteString(plug.Name)
		buf.WriteString("_out=")

		var opts []string
		for k, v := range plug.Options {
			if v != "" {
				opts = append(opts, fmt.Sprintf("%s=%s", k, v))
			} else {
				opts = append(opts, k)
			}
		}

		if len(opts) > 0 {
			buf.WriteString(strings.Join(opts, ","))
			buf.WriteByte(':')
		}

		buf.WriteString(plug.Out)
	}

	for _, file := range q.Files {
		buf.WriteByte(' ')
		buf.WriteString(file)
	}

	return buf.String()
}

type Core struct {
	deps            []string
	logger          *zerolog.Logger
	plugins         []Plugin
	inputs          Inputs
	console         console.Console
	storage         storage.IStorage
	moduleConfig    moduleconfig.IModuleConfig
	lockFile        lockfile.ILockFile
	protoRoot       string
	generateOutDirs bool
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
	generateOutDirs bool,
) *Core {
	return &Core{
		deps:            deps,
		logger:          logger,
		plugins:         plugins,
		inputs:          inputs,
		console:         console,
		storage:         storage,
		moduleConfig:    moduleConfig,
		lockFile:        lockFile,
		protoRoot:       protoRoot,
		generateOutDirs: generateOutDirs,
	}
}

func (c *Core) Generate(ctx context.Context, root, directory string) error {
	q := Query{
		Compiler: "protoc",
		Imports: []string{
			toolbox.Coalesce(c.protoRoot, root),
		},
		Plugins: c.plugins,
	}

	for _, dep := range c.deps {
		modulePaths, err := c.getModulePath(ctx, dep)
		if err != nil {
			return rerrors.Wrap(err, "c.getModulePath")
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

func (c *Core) getModulePath(_ context.Context, requestedDependency string) (string, error) {
	module := models.NewModule(requestedDependency)

	isInstalled, err := c.storage.IsModuleInstalled(module)
	if err != nil {
		return "", rerrors.Wrap(err, "h.storage.IsModuleInstalled")
	}

	if !isInstalled {
		return "", fmt.Errorf("module %s is not installed. Please run `moti install` first", module.Name)
	}

	lockFileInfo, err := c.lockFile.Read(module.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "lockFile.Read")
	}

	installedPath := c.storage.GetInstallDir(module.Name, lockFileInfo.Version)

	return installedPath, nil
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
