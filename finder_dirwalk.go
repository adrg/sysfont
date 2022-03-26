//go:build go1.16
// +build go1.16

package sysfont

import (
	"io/fs"
	"path/filepath"
)

func findFonts(opts *FinderOpts, fn func(filename string) error) {
	walker := func(filename string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		return fn(filename)
	}

	// Traverse OS font directories.
	for _, dir := range opts.SearchPaths {
		if err := filepath.WalkDir(dir, walker); err != nil {
			continue
		}
	}
}
