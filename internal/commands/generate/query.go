package generate

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"go.redsock.ru/moti/internal/config"
)

type ProtocQuery struct {
	Files   []string
	Imports []string
	Plugins []config.Plugin
}

func (q ProtocQuery) Build() (command string, args []string) {
	command = "protoc"

	for _, imp := range slices.Sorted(maps.Keys(toUniqueMap(q.Imports))) {
		args = append(args, "-I "+imp)
	}

	for _, plug := range q.Plugins {
		arg := "--" + plug.Name + "_out="

		var opts []string
		for k, v := range plug.Opts {
			if v != "" {
				opts = append(opts, fmt.Sprintf("%s=%s", k, v))
			} else {
				opts = append(opts, k)
			}
		}

		if len(opts) > 0 {
			arg += strings.Join(opts, ",") + ":"
		}

		arg += plug.Out
		args = append(args, arg)
	}

	uniqueProtoFileDirs := make(map[string]struct{})
	for _, file := range q.Files {
		uniqueProtoFileDirs[filepath.Dir(file)] = struct{}{}
	}

	for file := range uniqueProtoFileDirs {
		args = append(args, file+"/*.proto")
	}

	return command, args
}
