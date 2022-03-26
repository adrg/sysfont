//go:build !go1.16
// +build !go1.16

package sysfont

import (
	"os"
	"path/filepath"
)

func findFonts(opts *FinderOpts, fn func(filename string) error) {
	walker := func(filename string, info os.FileInfo, err error) error {
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
		if err := filepath.Walk(dir, walker); err != nil {
			continue
		}
	}
}
