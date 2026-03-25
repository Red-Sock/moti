package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/commands"
	"go.redsock.ru/moti/internal/commands/generate"
	"go.redsock.ru/moti/internal/commands/install"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "moti",
		Short:   "moti - usage info",
		Long:    "moti - description info",
		Version: version.System(),
	}

	rootCmd.PersistentFlags().
		String(flags.Config, flags.DefaultConfigFilePath,
			"Specify the absolute or relative path to the configuration file for setting up the application.")

	addCommands(rootCmd,
		install.Command{},
		generate.Command{},
	)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("rootCmd.Execute")
	}
}

func addCommands(root *cobra.Command, handlers ...commands.Handler) {
	for _, h := range handlers {
		root.AddCommand(h.Command())
	}
}
