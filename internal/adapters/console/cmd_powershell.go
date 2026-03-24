//go:build windows

package console

import (
	"bytes"
	"context"
	"os/exec"
)

type powershell struct{}

func New() *powershell {
	return &powershell{}
}

// RunCmd executes a shell command.
func (powershell) RunCmd(ctx context.Context, dir string, command string, commandParams ...string) (string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.CommandContext(ctx, command, commandParams...)
	cmd.Dir = dir
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", &RunError{
			Command:       command,
			CommandParams: commandParams,
			Dir:           dir,
			Err:           err,
			Stderr:        stderr.String(),
		}
	}

	return stdout.String(), nil
}
