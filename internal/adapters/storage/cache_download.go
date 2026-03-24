package storage

import (
	"fmt"
	"os"

	"go.redsock.ru/moti/internal/core/models"
)

// CacheDownload create path to downloaded cache.
// Like $GOPATH/pkg/mod/cache/download
func (s *Storage) CreateCacheDownloadDir(cacheDownloadPaths models.CacheDownloadPaths) error {
	if err := os.MkdirAll(cacheDownloadPaths.CacheDownloadDir, dirPerm); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	return nil
}
