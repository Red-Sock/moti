package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/core/models"
)

type Core struct {
	commands.Env
}

func (c *Core) Install(ctx context.Context) error {
	for _, dep := range c.MotiConfig.Deps {
		module := models.NewModule(dep)
		err := c.InstallPackage(ctx, module)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) InstallPackage(ctx context.Context, requestedModule models.Module) error {
	isInstalled, err := c.Storage.IsModuleInstalled(requestedModule)
	if err != nil {
		return fmt.Errorf("c.storage.IsModuleInstalled: %w", err)
	}
	if isInstalled {
		return nil
	}

	cacheRepositoryDir, err := c.Storage.CreateCacheRepositoryDir(requestedModule.Name)
	if err != nil {
		return fmt.Errorf("c.storage.CreateCacheRepositoryDir: %w", err)
	}

	repo, err := git.New(ctx, requestedModule.Name, cacheRepositoryDir, c.Console)
	if err != nil {
		return fmt.Errorf("git.New: %w", err)
	}

	revision, err := repo.ReadRevision(ctx, requestedModule.Version)
	if err != nil {
		return fmt.Errorf("repository.ReadRevision: %w", err)
	}

	cacheDownloadPaths := c.Storage.GetCacheDownloadPaths(requestedModule, revision)

	err = c.Storage.CreateCacheDownloadDir(cacheDownloadPaths)
	if err != nil {
		return fmt.Errorf("c.storage.CreateCacheDownloadDir: %w", err)
	}

	err = repo.Fetch(ctx, revision)
	if err != nil {
		return fmt.Errorf("repository.Fetch: %w", err)
	}

	moduleConfig, err := c.ModuleConfig.ReadFromRepo(ctx, repo, revision)
	if err != nil {
		return fmt.Errorf("c.moduleConfig.Read: %w", err)
	}

	for _, indirectDep := range moduleConfig.Dependencies {
		err = c.InstallPackage(ctx, indirectDep)
		if err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().Interface("dependency", indirectDep).Msg("Version not found")
				return models.ErrVersionNotFound
			}

			return fmt.Errorf("c.Get: %w", err)
		}
	}

	err = repo.Archive(ctx, revision, cacheDownloadPaths)
	if err != nil {
		return fmt.Errorf("repository.Archive: %w", err)
	}

	moduleHash, err := c.Storage.Install(cacheDownloadPaths, requestedModule, revision, moduleConfig)
	if err != nil {
		return fmt.Errorf("c.storage.Install: %w", err)
	}

	log.Debug().Str("hash", string(moduleHash)).Msg("HASH")

	err = c.LockFile.Write(requestedModule.Name, revision.Version, moduleHash)
	if err != nil {
		return fmt.Errorf("c.lockFile.Write: %w", err)
	}

	return nil
}
