package config

import (
	"go.redsock.ru/rerrors"
)

var ErrFileNotFound = rerrors.New("config file not found")

type Generate struct {
	Inputs    []Input  `json:"inputs" yaml:"inputs"`
	Plugins   []Plugin `json:"plugins" yaml:"plugins"`
	ProtoRoot string   `json:"proto_root" yaml:"proto_root"`
}

type Plugin struct {
	Name string            `json:"name" yaml:"name"`
	Out  string            `json:"out" yaml:"out"`
	Opts map[string]string `json:"opts" yaml:"opts"`
}

type Input struct {
	Directory string       `yaml:"directory"`
	GitRepo   InputGitRepo `yaml:"git_repo"`
}
type InputGitRepo struct {
	URL          string `yaml:"url"`
	SubDirectory string `yaml:"sub_directory"`
	Out          string `yaml:"out"`
}
