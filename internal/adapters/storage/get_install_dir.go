package storage

import (
	"path"

	"go.redsock.ru/moti/internal/helpers"
)

func (s *Storage) GetInstallDir(moduleName string, version string) string {
	version = helpers.SanitizePath(version)

	return path.Join(s.rootDir, installedDir, moduleName, version)
}
