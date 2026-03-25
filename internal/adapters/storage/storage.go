package storage

import (
	"context"

	"go.redsock.ru/moti/internal/models"
)

const (
	DirPerm = 0755
	// root cache dir
	cacheDir = "cache"
	// dir for downloaded (check sum, archive)
	cacheDownloadDir = "download"
	// dir for installed packages
	installedDir = "mod"
)

type LockFile interface {
	Read(moduleName string) (models.LockFileInfo, error)
}

// Storage implements workflows with directories
type Storage struct {
	rootDir  string
	lockFile LockFile
}

type IStorage interface {
	CreateCacheRepositoryDir(name string) (string, error)
	CreateCacheDownloadDir(models.CacheDownloadPaths) error
	GetCacheDownloadPaths(module models.Module, revision models.Revision) models.CacheDownloadPaths
	Install(
		ctx context.Context,
		cacheDownloadPaths models.CacheDownloadPaths,
		module models.Module,
		revision models.Revision,
		moduleConfig models.ModuleConfig,
	) (models.ModuleHash, error)
	GetInstalledModuleHash(moduleName string, revisionVersion string) (models.ModuleHash, error)
	IsModuleInstalled(module models.Module) (bool, error)
	GetInstallDir(moduleName string, revisionVersion string) string
}

func New(rootDir string, lockFile LockFile) *Storage {
	return &Storage{
		rootDir:  rootDir,
		lockFile: lockFile,
	}
}
