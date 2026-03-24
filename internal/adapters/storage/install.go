package storage

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"

	"github.com/codeclysm/extract/v3"
	"golang.org/x/mod/sumdb/dirhash"

	"go.redsock.ru/moti/internal/core/adapters"
	"go.redsock.ru/moti/internal/core/models"
)

// Install package from archive
// and calculateds hash of installed package
func (s *Storage) Install(
	cacheDownloadPaths models.CacheDownloadPaths,
	module models.Module,
	revision models.Revision,
	moduleConfig models.ModuleConfig,
) (models.ModuleHash, error) {
	log.Info().
		Str("package", module.Name).
		Str("version", revision.Version).
		Str("commit", revision.CommitHash).
		Msg("Install package")

	version := adapters.SanitizePath(revision.Version)
	installedDirPath := s.GetInstallDir(module.Name, version)

	if err := os.MkdirAll(installedDirPath, dirPerm); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	fp, err := os.Open(cacheDownloadPaths.ArchiveFile)
	if err != nil {
		return "", fmt.Errorf("os.Open: %w", err)
	}
	defer func() { _ = fp.Close() }()

	renamer := getRenamer(moduleConfig)

	log.Debug().Str("installedDirPath", installedDirPath).Msg("Starting extract")

	if err := extract.Archive(context.TODO(), fp, installedDirPath, renamer); err != nil {
		return "", fmt.Errorf("extract.Archive: %w", err)
	}

	installedPackageHash, err := dirhash.HashDir(installedDirPath, "", dirhash.DefaultHash)
	if err != nil {
		return "", fmt.Errorf("dirhash.HashDir: %w", err)
	}

	return models.ModuleHash(installedPackageHash), nil
}

// getRenamer return renamer function to convert result files path
func getRenamer(moduleConfig models.ModuleConfig) func(string) string {
	return func(file string) string {
		for _, dir := range moduleConfig.Directories {
			dir := dir + "/" // add trailing slash

			if strings.HasPrefix(file, dir) {
				return strings.TrimPrefix(file, dir)
			}
		}
		return file
	}
}
