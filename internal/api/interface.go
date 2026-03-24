package api

import (
	"github.com/spf13/cobra"
)

// Handler is an interface for a handling command.
type Handler interface {
	// Command returns a command.
	Command() *cobra.Command
}
