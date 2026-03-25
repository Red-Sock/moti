package storage

import (
	"os"

	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/models"
)

func (s *Storage) CreateCacheDownloadDir(cacheDownloadPaths models.CacheDownloadPaths) error {
	err := os.MkdirAll(cacheDownloadPaths.CacheDownloadDir, DirPerm)
	if err != nil {
		return rerrors.Wrap(err, "os.MkdirAll")
	}

	return nil
}
