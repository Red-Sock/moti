package models

import (
	"go.redsock.ru/rerrors"
)

type LockFileInfo struct {
	Name    string
	Version string
	Hash    ModuleHash
}

var ErrModuleNotFoundInLockFile = rerrors.New("module not found in lock file")
