package lockfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/adapters/fs"
	"go.redsock.ru/moti/internal/config"
	"go.redsock.ru/moti/internal/models"
)

type DirWalker interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	WalkDir(callback func(path string, err error) error) error
}

type fileInfo struct {
	version string
	hash    string
}

type LockFile struct {
	dirWalker DirWalker
	cache     map[string]fileInfo
}

//go:generate minimock -i ILockFile -o ../../mocks -g -s "_mock.go"
type ILockFile interface {
	Read(moduleName string) (models.LockFileInfo, error)
	Write(
		moduleName string, revisionVersion string, installedPackageHash models.ModuleHash,
	) error
}

func New(dirWalker DirWalker) (*LockFile, error) {
	cache := make(map[string]fileInfo)
	lockFile := &LockFile{
		dirWalker: dirWalker,
		cache:     cache,
	}

	lockFileOpened, err := dirWalker.Open(config.LockFileName)
	if err != nil {
		if !rerrors.Is(err, os.ErrNotExist) {
			return nil, rerrors.Wrap(err)
		}

		return lockFile, nil
	}

	fscanner := bufio.NewScanner(lockFileOpened)

	for fscanner.Scan() {
		parts := strings.Fields(fscanner.Text())
		if len(parts) != 3 {
			continue
		}

		fileInfo := fileInfo{
			version: parts[1],
			hash:    parts[2],
		}
		cache[parts[0]] = fileInfo
	}

	return lockFile, nil
}

func NewOrDie(workdir string) *LockFile {
	dirWalker := fs.NewFSWalker(workdir, ".")

	lock, err := New(dirWalker)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("could not create lock file")
	}

	return lock
}

// Read information about the module by its name from a lock file
func (l *LockFile) Read(moduleName string) (models.LockFileInfo, error) {
	fileInf, ok := l.cache[moduleName]
	if !ok {
		return models.LockFileInfo{}, models.ErrModuleNotFoundInLockFile
	}

	lockFileInfo := models.LockFileInfo{
		Name:    moduleName,
		Version: fileInf.version,
		Hash:    models.ModuleHash(fileInf.hash),
	}

	return lockFileInfo, nil
}

func (l *LockFile) Write(moduleName string, revisionVersion string, installedPackageHash models.ModuleHash) error {
	lockFile, err := l.dirWalker.Create(config.LockFileName)
	if err != nil {
		return fmt.Errorf("l.dirWalker.Create: %w", err)
	}

	fileInf := fileInfo{
		version: revisionVersion,
		hash:    string(installedPackageHash),
	}

	l.cache[moduleName] = fileInf

	keys := make([]string, 0, len(l.cache))
	for k := range l.cache {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		r := fmt.Sprintf("%s %s %s\n", k, l.cache[k].version, l.cache[k].hash)
		_, _ = lockFile.Write([]byte(r))
	}

	return nil
}
