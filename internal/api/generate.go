package api

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/fs/fs"
)

var _ Handler = (*Generate)(nil)

// deprecated: use internal/commands/generate.Command instead
type Generate struct{}

// deprecated: use internal/commands/generate.Command instead
func (g Generate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "generate code from proto files",
		Long:    "generate code from proto files",
		RunE:    g.Action,
	}

	cmd.Flags().StringP("path", "p", ".", "set path to directory with proto files")
	_ = cmd.MarkFlagRequired("path")

	return cmd
}

// deprecated: use internal/commands/generate.Command instead
func (g Generate) Action(cmd *cobra.Command, args []string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	dirWalker := fs.NewFSWalker(workingDir, ".")

	app, err := buildCore(cmd.Context(), *cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	dir, _ := cmd.Flags().GetString("path")
	// TODO somewhere here should be dependencies validation
	err = app.Generate(cmd.Context(), ".", dir)
	if err != nil {
		return fmt.Errorf("generator.Generate: %w", err)
	}

	return nil
}
