package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"

	"go.redsock.ru/moti/internal/api"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "protopack",
		Short:   "protopack - usage info",
		Long:    "protopack - description info",
		Version: version.System(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			isDebug, _ := cmd.Flags().GetBool(flags.DebugMode)
			initLogger(isDebug)
		},
	}

	rootCmd.PersistentFlags().String(flags.Config, flags.DefaultConfigFilePath, "Specify the absolute or relative path to the configuration file for setting up the application.")
	rootCmd.PersistentFlags().BoolP(flags.DebugMode, "d", false, "Enable debug mode to get more detailed information in logs.")

	addCommands(rootCmd,
		api.Lint{},
		api.Mod{},
		api.Completion{},
		api.Init{},
		api.Generate{},
		api.BreakingCheck{},
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("rootCmd.Execute")
	}
}

func initLogger(isDebug bool) {
	// use info as default level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func addCommands(root *cobra.Command, handlers ...api.Handler) {
	for _, h := range handlers {
		root.AddCommand(h.Command())
	}
}
