package test

import (
	"bytes"
	"io/fs"
)

func DirsEqual(dir1, dir2 fs.FS) (bool, error) {
	files1, err := getFiles(dir1)

	if err != nil {
		return false, err
	}

	files2, err := getFiles(dir2)

	if err != nil {
		return false, err
	}

	if len(files1) != len(files2) {
		return false, nil
	}

	for fname, payload1 := range files1 {
		payload2, ok := files2[fname]

		if !ok {
			return false, nil
		}

		if !bytes.Equal(payload1, payload2) {
			return false, nil
		}
	}

	return true, nil
}

// GetFiles reads the folder structure into the memory.
// This is just a test helper, no effort has been spent to optimize memory usage
func getFiles(dir fs.FS) (map[string][]byte, error) {
	out := map[string][]byte{}

	err := fs.WalkDir(dir, ".", func(p string, d fs.DirEntry, err2 error) error {
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
