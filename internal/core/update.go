package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"

	"go.redsock.ru/moti/internal/core/models"
)

// Update all packages from config
// dependencies slice of strings format: origin@version: github.com/company/repository@v1.2.3
// if version is absent use the latest commit
func (c *Core) Update(ctx context.Context, dependencies []string) error {
	for _, dependency := range dependencies {

		module := models.NewModule(dependency)

		isInstalled, err := c.storage.IsModuleInstalled(module)
		if err != nil {
			return fmt.Errorf("c.isModuleInstalled: %w", err)
		}

		if isInstalled {
			log.Info().Str("name", module.Name).Str("version", string(module.Version)).Msg("Module is installed")
			continue
		}

		if err := c.InstallPackage(ctx, module); err != nil {
			if errors.Is(err, models.ErrVersionNotFound) {
				log.Error().Str("dependency", dependency).Msg("Version not found")
				return models.ErrVersionNotFound
			}

			return fmt.Errorf("c.Get: %w", err)
		}
	}

	return nil
}
