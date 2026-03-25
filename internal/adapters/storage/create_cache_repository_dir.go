package storage

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

func (s *Storage) CreateCacheRepositoryDir(name string) (string, error) {
	cacheDirPath := filepath.Join(s.rootDir, cacheDir, fmt.Sprintf("%x", sha256.Sum256([]byte(name))))

	err := os.MkdirAll(cacheDirPath, DirPerm)
	if err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	return cacheDirPath, nil
}
