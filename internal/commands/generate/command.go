package generate

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/adapters/fs"
	"go.redsock.ru/moti/internal/commands"
)

type Command struct{}

func (c Command) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "generate code from proto files",
		Long:    "generate code from proto files",
		RunE:    c.Action,
	}

	cmd.Flags().StringP("path", "p", ".", "set path to directory with proto files")

	return cmd
}

func (c Command) Action(cmd *cobra.Command, args []string) error {
	err := c.Do(cmd, args)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to generate")
	}

	return nil
}

func (c Command) Do(cmd *cobra.Command, _ []string) error {
	app := Core{
		Env:    commands.GetProductionEnvironmentOrDie(cmd),
		Walker: &fs.FsWalker{},
	}

	err := app.Generate(cmd.Context())
	if err != nil {
		return fmt.Errorf("generator.Generate: %w", err)
	}

	return nil
}
