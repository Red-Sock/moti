package commands

import (
	"github.com/spf13/cobra"

	"go.redsock.ru/moti/internal/adapters/console"
	"go.redsock.ru/moti/internal/adapters/fs"
	lockfile "go.redsock.ru/moti/internal/adapters/lock_file"
	moduleconfig "go.redsock.ru/moti/internal/adapters/module_config"
	"go.redsock.ru/moti/internal/adapters/storage"
	"go.redsock.ru/moti/internal/config"
)

type Env struct {
	WorkDir    string
	MotiConfig config.Config

	Console console.Console

	Storage      storage.IStorage
	ModuleConfig moduleconfig.IModuleConfig
	LockFile     lockfile.ILockFile
}

func GetProductionEnvironmentOrDie(cmd *cobra.Command) (e Env) {
	e.Console = console.New()
	e.ModuleConfig = moduleconfig.New()

	e.WorkDir = fs.GetWdOrDie()
	e.MotiConfig = config.ReadOrDie(cmd)
	e.LockFile = lockfile.NewOrDie(e.WorkDir)
	e.Storage = storage.New(e.MotiConfig.CachePath, e.LockFile)

	return e
}
