package console

import (
	"context"
)

type Console interface {
	RunCmd(ctx context.Context, dir string, command string, commandParams ...string) (string, error)
}
