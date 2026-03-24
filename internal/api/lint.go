package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/core"
	"go.redsock.ru/moti/internal/flags"
	"go.redsock.ru/moti/internal/fs/fs"
)

var _ Handler = (*Lint)(nil)

// Lint is a handler for lint command.
type Lint struct{}

// Format is the format of output.
const (
	TextFormat = "text"
	JSONFormat = "json"
)

var (
	ErrHasLintIssue = errors.New("has lint issue")
)

// Command implements Handler.
func (l Lint) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lint",
		Aliases: []string{"l"},
		Short:   "linting proto files",
		Long:    "linting proto files",
		RunE:    l.Action,
	}

	cmd.Flags().StringP("path", "p", ".", "set relative path to directory with proto files")
	_ = cmd.MarkFlagRequired("path")

	cmd.Flags().StringP("format", "f", TextFormat, "set format of output (text, json)")

	return cmd
}

// Action implements Handler.
func (l Lint) Action(cmd *cobra.Command, args []string) error {
	err := l.action(cmd)
	if err != nil {
		var e *core.OpenImportFileError

		switch {
		case errors.Is(err, ErrHasLintIssue):
			os.Exit(1)
		case errors.As(err, &e):
			errExit(2, "Cannot import file", "file name", e.FileName)
		default:
			return err
		}
	}

	return nil
}

func (l Lint) action(cmd *cobra.Command) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	path, _ := cmd.Flags().GetString("path")
	configPath, _ := cmd.Flags().GetString(flags.Config)

	cfg, err := config.New(cmd.Context(), configPath)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}
	core.SetAllowCommentIgnores(cfg.Lint.AllowCommentIgnores)

	fsWalker := fs.NewFSWalker(workingDir, path)

	app, err := buildCore(cmd.Context(), *cfg, fsWalker)
	if err != nil {
		return fmt.Errorf("buildCore: %w", err)
	}
	issues, err := app.Lint(cmd.Context(), fsWalker)
	if err != nil {
		return fmt.Errorf("c.Lint: %w", err)
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

	return ErrHasLintIssue
}

func printIssues(format string, w io.Writer, issues []core.IssueInfo) error {
	switch format {
	case TextFormat:
		return textPrinter(w, issues)
	case JSONFormat:
		return jsonPrinter(w, issues)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// textPrinter prints the error in text format.
func textPrinter(w io.Writer, issues []core.IssueInfo) error {
	buffer := bytes.NewBuffer(nil)
	for _, issue := range issues {
		buffer.Reset()

		_, _ = buffer.WriteString(fmt.Sprintf("%s:%d:%d:%s %s (%s)",
			issue.Path,
			issue.Position.Line,
			issue.Position.Column,
			issue.SourceName,
			issue.Message,
			issue.RuleName,
		))
		_, _ = buffer.WriteString("\n")
		if _, err := w.Write(buffer.Bytes()); err != nil {
			return fmt.Errorf("w.Write: %w", err)
		}
	}

	return nil
}

// jsonPrinter prints the error in json format.
func jsonPrinter(w io.Writer, issues []core.IssueInfo) error {
	for _, issue := range issues {
		marshalErr := json.NewEncoder(w).Encode(issue)
		if marshalErr != nil {
			return fmt.Errorf("json.NewEncoder.Encode: %w", marshalErr)
		}
	}

	return nil
}
