package core

import (
	"context"
	"fmt"

	cp "github.com/otiai10/copy"
)

const (
	// TODO: move to config
	protopackVendorDir = "protopack_vendor"
)

// Vendor copy all proto files from deps to local dir
func (c *Core) Vendor(ctx context.Context) error {
	if err := c.Download(ctx, c.deps); err != nil {
		return fmt.Errorf("c.Download: %w", err)
	}

	for dep := range c.lockFile.DepsIter() {
		depPath, err := c.getModulePath(ctx, dep.Name)
		if err != nil {
			return fmt.Errorf("c.moduleReflect.GetModulePath: %w", err)
		}

		if err := cp.Copy(depPath, protopackVendorDir); err != nil {
			return fmt.Errorf("c.Copy: %w", err)
		}
	}

	return nil
}
