package generate

import (
	"fmt"
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
	q.Imports = removeDoubles(q.Imports)

	for _, imp := range q.Imports {
		//slices.Sorted(maps.Keys(toUniqueMap(q.Imports)))
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

	args = append(args, q.Files...)
	//uniqueProtoFileDirs := make(map[string]struct{})
	//for _, file := range q.Files {
	//	uniqueProtoFileDirs[filepath.Dir(file)] = struct{}{}
	//}
	//
	//for file := range uniqueProtoFileDirs {
	//	args = append(args, file+"/*.proto")
	//}

	return command, args
}

func removeDoubles(in []string) []string {
	out := make([]string, 0, len(in))
	existing := map[string]struct{}{}

	for _, file := range in {
		_, exists := existing[file]
		if exists {
			continue
		}
		out = append(out, file)
		existing[file] = struct{}{}
	}
	return out
}
