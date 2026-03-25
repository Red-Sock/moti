package install

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/models"
)

type Core struct {
	commands.Env
}

func (c *Core) Install(ctx context.Context) error {
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

	cacheRepositoryDir, err := c.Storage.CreateCacheRepositoryDir(requestedModule.Name)
	if err != nil {
		return rerrors.Wrap(err, "c.storage.CreateCacheRepositoryDir")
	}

	repo, err := git.New(ctx, requestedModule.Name, cacheRepositoryDir, c.Console)
	if err != nil {
		return rerrors.Wrap(err, "git.New: %w", err)
	}

	revision, err := repo.ReadRevision(ctx, requestedModule.Version)
	if err != nil {
		return rerrors.Wrap(err, "repository.ReadRevision")
	}

	cacheDownloadPaths := c.Storage.GetCacheDownloadPaths(requestedModule, revision)

	err = c.Storage.CreateCacheDownloadDir(cacheDownloadPaths)
	if err != nil {
		return rerrors.Wrap(err, "c.storage.CreateCacheDownloadDir")
	}

	err = repo.Fetch(ctx, revision)
	if err != nil {
		return rerrors.Wrap(err, "repository.Fetch")
	}

	moduleConfig, err := c.ModuleConfig.ReadFromRepo(ctx, repo, revision)
	if err != nil {
		return rerrors.Wrap(err, "c.moduleConfig.Read")
	}

	for _, indirectDep := range moduleConfig.Dependencies {
		err = c.InstallPackage(ctx, indirectDep)
		if err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().
					Interface("dependency", indirectDep).
					Msg("Version not found")

				return models.ErrVersionNotFound
			}

			return rerrors.Wrap(err, "c.Get")
		}
	}

	err = repo.Archive(ctx, revision, cacheDownloadPaths)
	if err != nil {
		return rerrors.Wrap(err, "repository.Archive")
	}

	moduleHash, err := c.Storage.Install(cacheDownloadPaths, requestedModule, revision, moduleConfig)
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
