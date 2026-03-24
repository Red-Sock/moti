// Package core contains every logic for working cli.
package core

import (
	"errors"
	"github.com/rs/zerolog"
)

// Core provide to business logic of ProtoPack.
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

	protoRoot       string
	generateOutDirs bool
}

var (
	ErrInvalidRule            = errors.New("invalid rule")
	ErrRepositoryDoesNotExist = errors.New("repository does not exist")
)

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
	generateOutDirs bool,
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
		generateOutDirs:         generateOutDirs,
	}
}
