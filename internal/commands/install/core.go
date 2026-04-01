package install

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.redsock.ru/moti/internal/adapters/console"
	"go.redsock.ru/moti/internal/adapters/repository"
	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/models"
)

type Core struct {
	commands.Env
	RepoFactory RepoFactory
}

//go:generate minimock -i RepoFactory -o ../../mocks -g -s "_mock.go"
type RepoFactory interface {
	New(ctx context.Context, remote string, cacheDir string, console git.Console) (repository.Repo, error)
}

type gitRepoFactory struct{}

func (f *gitRepoFactory) New(ctx context.Context, remote string, cacheDir string, console git.Console) (repository.Repo, error) {
	return git.New(ctx, remote, cacheDir, console)
}

func (c *Core) Install(ctx context.Context) error {
	if c.RepoFactory == nil {
		c.RepoFactory = &gitRepoFactory{}
	}

	for _, installBin := range c.MotiConfig.Binaries.Install {
		if installBin.Go.Module != "" {
			module := models.NewModule(installBin.Go.Module)
			version := toolbox.Coalesce(string(module.Version), "latest")

			isInstalled, err := c.isGoBinaryVersionInstalled(ctx, installBin.Go, module, version)
			if err != nil {
				return rerrors.Wrap(err, "isGoBinaryVersionInstalled")
			}

			if isInstalled {
				log.Info().
					Str("module", installBin.Go.Module).
					Str("version", version).
					Msg("module already installed")
				continue
			}

			log.Info().
				Str("module", installBin.Go.Module).
				Str("version", version).
				Msg("module is not installed. Installing...")

			err = c.installGoBin(ctx, installBin.Go, module, version)
			if err != nil {
				return rerrors.Wrap(err, "error installing go binary")
			}
		}
	}

	for _, dep := range c.MotiConfig.Deps {
		module := models.NewModule(dep)

		err := c.InstallPackage(ctx, module)
		if err != nil {
			return rerrors.Wrap(err)
		}
	}

	return nil
}

func (c *Core) InstallPackage(ctx context.Context, requestedModule models.Module) error {
	isInstalled, err := c.Storage.IsModuleInstalled(requestedModule)
	if err != nil {
		return rerrors.Wrap(err, "c.storage.IsModuleInstalled")
	}

	if isInstalled {
		return nil
	}

	repo, revision, err := c.fetchAndReadRevision(ctx, requestedModule)
	if err != nil {
		return rerrors.Wrap(err, "fetchAndReadRevision")
	}

	moduleConfig, err := c.ModuleConfig.ReadFromRepo(ctx, repo, revision)
	if err != nil {
		return rerrors.Wrap(err, "c.moduleConfig.Read")
	}

	err = c.installDependencies(ctx, moduleConfig.Dependencies)
	if err != nil {
		return rerrors.Wrap(err, "installDependencies")
	}

	cacheDownloadPaths := c.Storage.GetCacheDownloadPaths(requestedModule, revision)

	err = c.Storage.CreateCacheDownloadDir(cacheDownloadPaths)
	if err != nil {
		return rerrors.Wrap(err, "c.storage.CreateCacheDownloadDir")
	}

	err = repo.Archive(ctx, revision, cacheDownloadPaths)
	if err != nil {
		return rerrors.Wrap(err, "repository.Archive")
	}

	moduleHash, err := c.Storage.Install(ctx, cacheDownloadPaths, requestedModule, revision, moduleConfig)
	if err != nil {
		return rerrors.Wrap(err, "c.storage.Install")
	}

	log.Debug().
		Str("hash", string(moduleHash)).
		Msg("Hash for module")

	err = c.LockFile.Write(requestedModule.Name, revision.Version, moduleHash)
	if err != nil {
		return rerrors.Wrap(err, "c.lockFile.Write")
	}

	return nil
}

func (c *Core) fetchAndReadRevision(ctx context.Context, requestedModule models.Module) (
	repository.Repo, models.Revision, error) {
	cacheRepositoryDir, err := c.Storage.CreateCacheRepositoryDir(requestedModule.Name)
	if err != nil {
		return nil, models.Revision{}, rerrors.Wrap(err, "c.storage.CreateCacheRepositoryDir")
	}

	repo, err := c.RepoFactory.New(ctx, requestedModule.Name, cacheRepositoryDir, c.Console)
	if err != nil {
		return nil, models.Revision{}, rerrors.Wrap(err, "git.New")
	}

	revision, err := repo.ReadRevision(ctx, requestedModule.Version)
	if err != nil {
		return nil, models.Revision{}, rerrors.Wrap(err, "repository.ReadRevision")
	}

	err = repo.Fetch(ctx, revision)
	if err != nil {
		return nil, models.Revision{}, rerrors.Wrap(err, "repository.Fetch")
	}

	return repo, revision, nil
}

func (c *Core) installGoBin(ctx context.Context, goBin config.GoBin, module models.Module, version string) error {
	log.Info().
		Str("module", goBin.Module).
		Msg("Installing go binary")

	command := "go install " + module.Name + "@" + version

	if c.MotiConfig.Binaries.BinDir != "" {
		command = "GOBIN=" + c.MotiConfig.Binaries.BinDir + " " + command
	}

	_, err := c.Console.RunCmd(ctx, c.WorkDir, command)
	if err != nil {
		return rerrors.Wrap(err, "error executing go install")
	}

	return nil
}

func (c *Core) isGoBinaryVersionInstalled(ctx context.Context,
	goBin config.GoBin,
	module models.Module,
	expectedVersion string,
) (bool, error) {

	gobin := c.MotiConfig.BuildGOBIN(c.Env.WorkDir)[len(config.GOBINPrefix):]
	binaryName := gobin + "/" + module.Name[strings.LastIndex(module.Name, "/")+1:]

	versionCheckArgs := toolbox.Coalesce(goBin.VersionCheckArgs, "")

	output, err := c.Console.RunCmd(ctx, c.WorkDir, binaryName, versionCheckArgs)
	if err != nil {
		if errors.Is(err, console.ErrNotFound) {
			return false, nil
		}
		return false, rerrors.Wrap(err, "error checking binary version")
	}

	if expectedVersion != "latest" && !strings.Contains(output, expectedVersion) {
		return false, rerrors.New("version check failed: output " + output + " does not contain version " + expectedVersion)
	}

	return true, nil
}

func (c *Core) installDependencies(ctx context.Context, dependencies []models.Module) error {
	for _, indirectDep := range dependencies {
		err := c.InstallPackage(ctx, indirectDep)
		if err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().
					Interface("dependency", indirectDep).
					Msg("Version not found")

				return models.ErrVersionNotFound
			}

			return rerrors.Wrap(err, "c.InstallPackage (indirect)")
		}
	}

	return nil
}
