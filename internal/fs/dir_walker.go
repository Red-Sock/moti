package fs

import (
	"io/fs"
	"os"

	"go.redsock.ru/rerrors"
)

type Walker struct {
	*Adapter

	path string
}

func NewFSWalker(root, path string) *Walker {
	if path == "" {
		path = "."
	}

	diskFS := os.DirFS(root)

	return &Walker{
		Adapter: &Adapter{diskFS, root},
		path:    path,
	}
}

func (w *Walker) WalkDir(callback func(path string, err error) error) error {
	err := fs.WalkDir(w.FS, w.path, func(path string, d fs.DirEntry, err error) error {
		return callback(path, err)
	})
	if err != nil {
		return rerrors.Wrap(err, "fs.WalkDir")
	}

	return nil
}
