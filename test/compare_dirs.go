package test

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"slices"
)

var ignoredFnames = map[string]struct{}{
	".":        {},
	".gitkeep": {},
}

func DirsEqual(dir1, dir2 fs.FS) error {
	files1, err := getFiles(dir1)

	if err != nil {
		return err
	}

	files2, err := getFiles(dir2)

	if err != nil {
		return err
	}

	if len(files1) != len(files2) {
		return fmt.Errorf("File lists do not match. Dir1: %v, Dir2: %v", keys(files1), keys(files2))
	}

	for fname, payload1 := range files1 {
		payload2, ok := files2[fname]

		if !ok {
			return fmt.Errorf("File [%s] is missing in dir2", fname)
		}

		if !bytes.Equal(payload1, payload2) {
			return fmt.Errorf("File [%s] payload is different: \ndir1:\n[%s]\n\ndir2:\n[%s]", fname, string(payload1), string(payload2))
		}
	}

	return nil
}

func keys[A any](m map[string]A) []string {
	out := make([]string, 0, len(m))

	for k := range m {
		out = append(out, k)
	}

	// for repeatable output
	slices.Sort(out)

	return out
}

// GetFiles reads the folder structure into the memory.
// This is just a test helper, no effort has been spent to optimize memory usage
func getFiles(dir fs.FS) (map[string][]byte, error) {
	out := map[string][]byte{}

	err := fs.WalkDir(dir, ".", func(p string, d fs.DirEntry, err2 error) error {
		if _, ok := ignoredFnames[path.Base(p)]; ok {
			return nil
		}

		if d.IsDir() {
			out[p] = nil
			return nil
		}

		if err2 != nil {
			return err2
		}

		content, err := fs.ReadFile(dir, p)

		if err != nil {
			return nil
		}

		out[p] = content
		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
