package generate

import (
	"errors"
	"strings"

	"go.redsock.ru/moti/internal/adapters/console"
)

type FileNotFoundOrHadErrors struct {
	desc string
}

func (f FileNotFoundOrHadErrors) Error() string {
	return `
|============================================|
|You received "File Not found" error 		 |
|Most likely you messed your imports. 		 |
|Read proper import technics at readme		 |
|           https://github.com/Red-Sock/moti |
|============================================|
Error itself:
` + f.desc
}

type MissingGoImportPathError struct {
	desc string
}

func (m MissingGoImportPathError) Error() string {
	return `
|================================|
| Missing Go package path.		 | 
| Pass it with go_package option |
|================================|
Error itself:
` + m.desc
}

func parseError(in error) error {
	var consoleErr *console.RunError
	if !errors.As(in, &consoleErr) {
		return in
	}

	if strings.Contains(consoleErr.Stderr, "File not found") {
		return FileNotFoundOrHadErrors{
			desc: consoleErr.Stderr,
		}
	}

	if strings.Contains(consoleErr.Stderr, "unable to determine Go import path for") {
		return MissingGoImportPathError{
			desc: consoleErr.Stderr,
		}
	}

	return in
}

//12:33PM ERR failed to generate error="generator.Generate: Command: PATH=/Users/alexbukov/redsock/moti/examples/full/bin:$PATH protoc; Err: exit status 1; Stderr: protoc-gen-go: unable to determine Go import path for \"messages.proto\"\n\nPlease specify either:\n\t• a \"go_package\" option in the .proto source file, or\n\t• a \"M\" argument on the command line.\n\nSee https://protobuf.dev/reference/go/go-generated#package for more information.\n\n--go_out: protoc-gen-go: Plugin failed with status code 1.\n\nadapters.RunCmd"
