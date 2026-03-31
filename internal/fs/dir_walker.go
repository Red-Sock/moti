package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.redsock.ru/rerrors"
)

type DirWalker struct {
	*Adapter

	path string
}

func NewFSWalker(root, path string) *DirWalker {
	if path == "" {
		path = "."
	}

	diskFS := os.DirFS(root)

	return &DirWalker{
		Adapter: &Adapter{diskFS, root},
		path:    path,
	}
}

func (w *DirWalker) WalkDir(callback func(path string, err error) error) error {
	err := fs.WalkDir(w.FS, w.path, func(path string, d fs.DirEntry, err error) error {
		return callback(path, err)
	})
	if err != nil {
		return rerrors.Wrap(err, "fs.WalkDir")
	}

	return nil
}

func (a *DirWalker) Create(name string) (io.WriteCloser, error) {
	path := filepath.Join(a.rootDir, name)

	f, err := os.Create(path)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	return f, nil
}
