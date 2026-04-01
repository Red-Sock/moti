package fs

import (
	"io"
	"io/fs"

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
