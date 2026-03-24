package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/core/models"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/fs/fs"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/core"
)

var _ Handler = (*Mod)(nil)

// Mod is a handler for package manager
type Mod struct{}

func (m Mod) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mod",
		Aliases: []string{"m"},
		Short:   "package manager",
		Long:    "package manager",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "download",
		Short: "download modules to local cache",
		RunE:  m.Download,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "update modules version using version from config",
		RunE:  m.Update,
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "vendor",
		Short: "copy proto files from deps to vendor dir",
		RunE:  m.Vendor,
	})

	return cmd
}

func (m Mod) Download(cmd *cobra.Command, args []string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}
	dirWalker := fs.NewFSWalker(workingDir, ".")

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}
	core.SetAllowCommentIgnores(cfg.Lint.AllowCommentIgnores)

	app, err := buildCore(cmd.Context(), *cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	if err := app.Download(cmd.Context(), cfg.Deps); err != nil {
		if errors.Is(err, models.ErrVersionNotFound) {
			os.Exit(1)
		}

		return fmt.Errorf("cmd.Download: %w", err)
	}
	return nil
}

func (m Mod) Update(cmd *cobra.Command, args []string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}
	dirWalker := fs.NewFSWalker(workingDir, ".")

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}
	core.SetAllowCommentIgnores(cfg.Lint.AllowCommentIgnores)

	app, err := buildCore(cmd.Context(), *cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	if err := app.Update(cmd.Context(), cfg.Deps); err != nil {
		if errors.Is(err, models.ErrVersionNotFound) {
			os.Exit(1)
		}

		return fmt.Errorf("cmd.Download: %w", err)
	}
	return nil
}

func (m Mod) Vendor(cmd *cobra.Command, args []string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}
	dirWalker := fs.NewFSWalker(workingDir, ".")

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}
	core.SetAllowCommentIgnores(cfg.Lint.AllowCommentIgnores)

	app, err := buildCore(cmd.Context(), *cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	if err := app.Vendor(cmd.Context()); err != nil {
		if errors.Is(err, models.ErrVersionNotFound) {
			os.Exit(1)
		}

		return fmt.Errorf("cmd.Download: %w", err)
	}
	return nil
}
