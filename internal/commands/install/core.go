package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.redsock.ru/moti/internal/adapters/console"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/repository/git"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/core/models"
)

type Core struct {
	logger       *zerolog.Logger
	console      console.Console
	storage      storage.IStorage
	moduleConfig moduleconfig.IModuleConfig
	lockFile     lockfile.ILockFile
}

func New(
	logger *zerolog.Logger,
	console console.Console,
	storage storage.IStorage,
	moduleConfig moduleconfig.IModuleConfig,
	lockFile lockfile.ILockFile,
) *Core {
	return &Core{
		logger:       logger,
		console:      console,
		storage:      storage,
		moduleConfig: moduleConfig,
		lockFile:     lockFile,
	}
}

func (c *Core) Install(ctx context.Context, deps []string) error {
	for _, dep := range deps {
		module := models.NewModule(dep)
		err := c.InstallPackage(ctx, module)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) InstallPackage(ctx context.Context, requestedModule models.Module) error {
	isInstalled, err := c.storage.IsModuleInstalled(requestedModule)
	if err != nil {
		return fmt.Errorf("c.storage.IsModuleInstalled: %w", err)
	}
	if isInstalled {
		return nil
	}

	cacheRepositoryDir, err := c.storage.CreateCacheRepositoryDir(requestedModule.Name)
	if err != nil {
		return fmt.Errorf("c.storage.CreateCacheRepositoryDir: %w", err)
	}

	repo, err := git.New(ctx, requestedModule.Name, cacheRepositoryDir, c.console)
	if err != nil {
		return fmt.Errorf("git.New: %w", err)
	}

	revision, err := repo.ReadRevision(ctx, requestedModule.Version)
	if err != nil {
		return fmt.Errorf("repository.ReadRevision: %w", err)
	}

	cacheDownloadPaths := c.storage.GetCacheDownloadPaths(requestedModule, revision)

	err = c.storage.CreateCacheDownloadDir(cacheDownloadPaths)
	if err != nil {
		return fmt.Errorf("c.storage.CreateCacheDownloadDir: %w", err)
	}

	err = repo.Fetch(ctx, revision)
	if err != nil {
		return fmt.Errorf("repository.Fetch: %w", err)
	}

	moduleConfig, err := c.moduleConfig.ReadFromRepo(ctx, repo, revision)
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

	moduleHash, err := s_install(c.storage, cacheDownloadPaths, requestedModule, revision, moduleConfig)
	if err != nil {
		return fmt.Errorf("c.storage.Install: %w", err)
	}

	log.Debug().Str("hash", string(moduleHash)).Msg("HASH")

	err = c.lockFile.Write(requestedModule.Name, revision.Version, moduleHash)
	if err != nil {
		return fmt.Errorf("c.lockFile.Write: %w", err)
	}

	return nil
}

func s_install(s storage.IStorage, cacheDownloadPaths models.CacheDownloadPaths, module models.Module, revision models.Revision, moduleConfig models.ModuleConfig) (models.ModuleHash, error) {
	return s.Install(cacheDownloadPaths, module, revision, moduleConfig)
}
