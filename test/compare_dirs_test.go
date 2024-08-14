package test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestDirsEqual(t *testing.T) {
	var ex = []struct {
		dir1  fstest.MapFS
		dir2  fstest.MapFS
		equal bool
	}{
		{
			dir1: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			equal: true,
		},
		{
			dir1: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`1234`),
				},
			},
			equal: false,
		},
		{
			dir1: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
				"file2": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`1234`),
				},
			},
			equal: false,
		},
		{
			dir1: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"file1": &fstest.MapFile{
					Data: []byte(`1234`),
				},
			},
			equal: false,
		},
		{
			dir1: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			equal: true,
		},
		{
			dir1: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: nil,
				},
			},
			equal: false,
		},
		{
			dir1: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
			},
			dir2: fstest.MapFS{
				"dir/file1": &fstest.MapFile{
					Data: []byte(`123`),
				},
				"dir/dir2/dir3": &fstest.MapFile{
					Data: nil,
				},
			},
			equal: false,
		},
	}

	for idx, ex := range ex {
		equal, err := DirsEqual(ex.dir1, ex.dir2)
		assert.NoErrorf(t, err, "[example %d]", idx+1)

		assert.Equalf(t, ex.equal, equal, "[example %d]", idx+1)
	}
}
