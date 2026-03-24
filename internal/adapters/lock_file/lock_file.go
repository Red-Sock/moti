package lockfile

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"sort"
	"strings"

	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/core"
	"go.redsock.ru/moti/internal/core/models"
)

const (
	lockFileName = "protopack.lock"
)

type fileInfo struct {
	version string
	hash    string
}

type LockFile struct {
	dirWalker core.DirWalker
	cache     map[string]fileInfo
}

func New(dirWalker core.DirWalker) (*LockFile, error) {
	cache := make(map[string]fileInfo)
	lockFile := &LockFile{
		dirWalker: dirWalker,
		cache:     cache,
	}

	fp, err := dirWalker.Open(lockFileName)
	if err != nil {
		if !rerrors.Is(err, os.ErrNotExist) {
			return nil, rerrors.Wrap(err)
		}

		return lockFile, nil
	}

	fscanner := bufio.NewScanner(fp)

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

// Read information about module by its name from lock file
// github.com/grpc-ecosystem/grpc-gateway v0.0.0-20240502030614-85850831b7bad2b8b60cb09783d8095176f22d98 h1:hRu1vxAH6CVNmz12mpqKue5HVBQP2neoaM/q2DLm0i4=
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

func (l *LockFile) Write(
	moduleName string, revisionVersion string, installedPackageHash models.ModuleHash,
) error {
	fp, err := l.dirWalker.Create(lockFileName)
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
		_, _ = fp.Write([]byte(r))
	}

	return nil
}

func (l *LockFile) DepsIter() iter.Seq[models.LockFileInfo] {
	return func(yield func(models.LockFileInfo) bool) {
		for moduleName, fileInf := range l.cache {
			lockFileInfo := models.LockFileInfo{
				Name:    moduleName,
				Version: fileInf.version,
				Hash:    models.ModuleHash(fileInf.hash),
			}
			if !yield(lockFileInfo) {
				return
			}
		}
	}
}

// IsEmpty check if lock file doesn't have any deps
func (l *LockFile) IsEmpty() bool {
	return len(l.cache) == 0
}
