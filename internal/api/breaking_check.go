package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/core"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/fs/fs"
)

var _ Handler = (*BreakingCheck)(nil)

// BreakingCheck is a handler for breaking command
type BreakingCheck struct{}

var (
	ErrBreakingCheckIssue = errors.New("has breaking check issue")
)

func (b BreakingCheck) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "breaking",
		Short: "api breaking check",
		Long:  "api breaking check",
		RunE:  b.Action,
	}

	cmd.Flags().StringP("path", "p", ".", "set relative path to directory with proto files")
	_ = cmd.MarkFlagRequired("path")

	cmd.Flags().StringP("format", "f", TextFormat, "set format of output (text, json)")
	cmd.Flags().String("against", "master", "set branch to compare with")
	_ = cmd.MarkFlagRequired("against")

	return cmd
}

func (b BreakingCheck) Action(cmd *cobra.Command, args []string) error {
	err := b.action(cmd)
	if err != nil {
		var e *core.OpenImportFileError
		var g *core.GitRefNotFoundError

		switch {
		case errors.Is(err, ErrBreakingCheckIssue):
			os.Exit(1)
		case errors.As(err, &e):
			errExit(2, "Cannot import file", "file name", e.FileName)
		case errors.As(err, &g):
			errExit(2, "Cannot find git ref", "ref", g.GitRef)
		case errors.Is(err, core.ErrRepositoryDoesNotExist):
			errExit(2, "Repository does not exist in current directory")
		default:
			return err
		}
	}

	return nil
}

func (b BreakingCheck) action(cmd *cobra.Command) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	configPath, _ := cmd.Flags().GetString(flags.Config)
	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	against, _ := cmd.Flags().GetString("against")
	if against != "" {
		cfg.BreakingCheck.AgainstGitRef = against
	}

	dirWalker := fs.NewFSWalker(workingDir, ".")
	app, err := buildCore(cmd.Context(), *cfg, dirWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}

	issues, err := app.BreakingCheck(cmd.Context(), workingDir, path)
	if err != nil {
		return fmt.Errorf("app.BreakingCheck: %w", err)
	}

	if len(issues) == 0 {
		return nil
	}

	format, _ := cmd.Flags().GetString("format")
	if err := printIssues(
		format,
		os.Stdout,
		issues,
	); err != nil {
		return fmt.Errorf("printLintErrors: %w", err)
	}

	return ErrBreakingCheckIssue
}
