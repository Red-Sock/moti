package config

import (
	"go.redsock.ru/rerrors"
)

var ErrFileNotFound = rerrors.New("config file not found")

type Generate struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Plugins []Plugin `json:"plugins" yaml:"plugins"`
}

type Binaries struct {
	BinDir      string `json:"bin_dir" yaml:"bin_dir"`
	AllowCustom bool   `json:"allow_custom" yaml:"allow_custom"`
	Install     []struct {
		Go GoBin `json:"go" yaml:"go"`
	} `json:"install" yaml:"install"`
}

type GoBin struct {
	Module           string `json:"module" yaml:"module"`
	VersionCheckArgs string `json:"version_check_args" yaml:"version_check_args"`
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
