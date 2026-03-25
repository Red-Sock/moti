package install

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/commands"
)

type Command struct{}

func (c Command) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "install dependencies",
		Long:    "install dependencies specified in moti.yaml",
		RunE:    c.Action,
	}

	return cmd
}

func (c Command) Action(cmd *cobra.Command, _ []string) error {
	err := c.Do(cmd)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to install dependencies")
	}

	return nil
}

func (c Command) Do(cmd *cobra.Command) error {
	app := Core{
		Env: commands.GetProductionEnvironmentOrDie(cmd),
	}

	err := app.Install(cmd.Context())
	if err != nil {
		return fmt.Errorf("install: %w", err)
	}

	return nil
}
