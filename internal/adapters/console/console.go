package console

import (
	"context"
)

//go:generate minimock -i Console -o ../../mocks -g -s "_mock.go"
type Console interface {
	RunCmd(ctx context.Context, dir string, command string, commandParams ...string) (string, error)
}
