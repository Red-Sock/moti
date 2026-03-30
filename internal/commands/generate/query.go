package generate

import (
	"fmt"
	"strings"

	"go.redsock.ru/moti/internal/config"
)

const (
	protocBin = "protoc"
)

type ProtocQuery struct {
	Imports []string
	Plugins []config.Plugin

	Files []string
}

func (q ProtocQuery) Build() (command string, args []string) {
	args = q.buildImports()
	args = append(args, q.buildPlugins()...)
	args = append(args, q.Files...)

	//uniqueProtoFileDirs := make(map[string]struct{})
	//for _, file := range q.Files {
	//	uniqueProtoFileDirs[filepath.Dir(file)] = struct{}{}
	//}
	//
	//for file := range uniqueProtoFileDirs {
	//	args = append(args, file+"/*.proto")
	//}

	return protocBin, args
}

func (q ProtocQuery) buildImports() (imports []string) {
	q.Imports = removeDoubles(q.Imports)

	imports = make([]string, 0, len(q.Imports))

	for _, imp := range q.Imports {
		imports = append(imports, "-I "+imp)
	}

	return imports
}

func (q ProtocQuery) buildPlugins() (plugins []string) {
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
		plugins = append(plugins, arg)
	}

	return plugins
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
