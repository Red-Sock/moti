package console

import (
	"fmt"
)

type RunError struct {
	Command       string
	CommandParams []string
	Dir           string
	Err           error
	Stderr        string
}

func (e RunError) Error() string {
	return fmt.Sprintf("Err: %v; Stderr: %s", e.Err, e.Stderr)
}
