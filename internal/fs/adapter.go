package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.redsock.ru/rerrors"
)

type Adapter struct {
	fs.FS

	rootDir string
}

func (a *Adapter) Open(name string) (io.ReadCloser, error) {
	rc, err := a.FS.Open(name)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	return rc, nil
}

func (a *Adapter) Create(name string) (io.WriteCloser, error) {
	path := filepath.Join(a.rootDir, name)

	f, err := os.Create(path)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	return f, nil
}
