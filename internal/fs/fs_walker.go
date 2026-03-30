package fs

//go:generate minimock -i IWalker -o ../../mocks -g -s "_mock.go"
type IWalker interface {
	WalkDir(root string, callback func(path string, err error) error) error
}

type FsWalker struct{}

func (f *FsWalker) WalkDir(root string, callback func(path string, err error) error) error {
	w := NewFSWalker(root, "")
	return w.WalkDir(callback)
}
