package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"

	"go.redsock.ru/moti/internal/core/models"
)

// Download all packages from config
// dependencies slice of strings format: origin@version: github.com/company/repository@v1.2.3
// if version is absent use the latest commit
func (c *Core) Download(ctx context.Context, dependencies []string) error {
	if c.lockFile.IsEmpty() {
		// if lock file is empty or doesn't exist install versions
		// from moti.yaml config and create lock file
		log.Debug().Msg("Lock file is empty")
		return c.Update(ctx, dependencies)
	}

	log.Debug().Msg("Lock file is not empty. Install deps from it")

	for lockFileInfo := range c.lockFile.DepsIter() {
		module := models.NewModuleFromLockFileInfo(lockFileInfo)

		isInstalled, err := c.storage.IsModuleInstalled(module)
		if err != nil {
			return fmt.Errorf("c.isModuleInstalled: %w", err)
		}

		if isInstalled {
			log.Info().Str("name", module.Name).Str("version", string(module.Version)).Msg("Module is installed")
			continue
		}

		err = c.InstallPackage(ctx, module)
		if err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().Str("name", module.Name).Str("version", string(module.Version)).Msg("Version not found")
				return models.ErrVersionNotFound
			}

			return fmt.Errorf("c.Get: %w", err)
		}
	}

	return nil
}
