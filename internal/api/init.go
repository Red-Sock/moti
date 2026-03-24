package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/fs/fs"
)

var _ Handler = (*Init)(nil)

// Init is a handler for initialization ProtoPack configuration.
type Init struct{}

// Command implements Handler.
func (i Init) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "initialize configuration",
		Long:    "initialize configuration",
		RunE:    i.Action,
	}

	cmd.Flags().StringP("dir", "d", ".", "directory path to initialize")
	_ = cmd.MarkFlagRequired("dir")

	return cmd
}

// Action implements Handler.
func (i Init) Action(cmd *cobra.Command, args []string) error {
	rootPath, _ := cmd.Flags().GetString("dir")
	dirFS := fs.NewFSWalker(rootPath, ".")

	cfg := &config.Config{}

	app, err := buildCore(cmd.Context(), *cfg, dirFS)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	err = app.Initialize(cmd.Context(), dirFS, []string{"DEFAULT"})
	if err != nil {
		return fmt.Errorf("initer.Initialize: %w", err)
	}

	return nil
}
