// Package core contains every logic for working cli.
package core

import (
	"errors"

	"github.com/rs/zerolog"
)

// deprecated: Use internal/commands/generate.Core instead for generation
type Core struct {
	rules        []Rule
	ignore       []string
	deps         []string
	ignoreOnly   map[string][]string
	logger       *zerolog.Logger
	plugins      []Plugin
	inputs       Inputs
	console      Console
	storage      Storage
	moduleConfig ModuleConfig
	lockFile     LockFile

	breakingCheckConfig     BreakingCheckConfig
	currentProjectGitWalker CurrentProjectGitWalker

	protoRoot string
}

var (
	ErrInvalidRule            = errors.New("invalid rule")
	ErrRepositoryDoesNotExist = errors.New("repository does not exist")
)

// deprecated: Use internal/commands/generate.New instead for generation
func New(
	rules []Rule,
	ignore []string,
	deps []string,
	ignoreOnly map[string][]string,
	logger *zerolog.Logger,
	plugins []Plugin,
	inputs Inputs,
	console Console,
	storage Storage,
	moduleConfig ModuleConfig,
	lockFile LockFile,
	currentProjectGitWalker CurrentProjectGitWalker,
	breakingCheckConfig BreakingCheckConfig,
	protoRoot string,
) *Core {
	return &Core{
		rules:                   rules,
		ignore:                  ignore,
		deps:                    deps,
		ignoreOnly:              ignoreOnly,
		logger:                  logger,
		plugins:                 plugins,
		inputs:                  inputs,
		console:                 console,
		storage:                 storage,
		moduleConfig:            moduleConfig,
		lockFile:                lockFile,
		currentProjectGitWalker: currentProjectGitWalker,
		breakingCheckConfig:     breakingCheckConfig,
		protoRoot:               protoRoot,
	}
}
