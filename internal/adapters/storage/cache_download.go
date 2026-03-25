package storage

import (
	"fmt"
	"os"

	"go.redsock.ru/moti/internal/models"
)

func (s *Storage) CreateCacheDownloadDir(cacheDownloadPaths models.CacheDownloadPaths) error {
	if err := os.MkdirAll(cacheDownloadPaths.CacheDownloadDir, DirPerm); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	return nil
}
