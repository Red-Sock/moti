//go:build unix

package console

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

type bash struct{}

func New() *bash {
	return &bash{}
}

func (bash) RunCmd(ctx context.Context, dir string, command string, commandParams ...string) (string, error) {
	var stderr, stdout bytes.Buffer

	fullCommand := append([]string{command}, commandParams...)

	cmd := exec.CommandContext(ctx, "bash", "-c", strings.Join(fullCommand, " "))
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
