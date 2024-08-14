package test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/can3p/go-hrx/hrx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const hrxPath = "../spec/example"
const extractedPath = "../spec/example/extracted"
const invalidPath = "../spec/example/invalid"

func getHrxFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, d := range files {
		if !strings.HasSuffix(d.Name(), ".hrx") {
			continue
		}

		names = append(names, d.Name())
	}

	return names, nil
}

func TestSpec(t *testing.T) {
	testFiles, err := getHrxFiles(hrxPath)
	require.NoError(t, err)

	t.Log(testFiles)
	assert.Positive(t, len(testFiles))

	for idx, fname := range testFiles {
		dirname, found := strings.CutSuffix(fname, ".hrx")
		require.Truef(t, found, "[example %d]", idx+1)
		resultFolder := path.Join(extractedPath, dirname)

		reader, err := hrx.OpenReader(fname)

		assert.NoErrorf(t, err, "[example %d]", idx+1)

		if err != nil {
			continue
		}

		equal, err := DirsEqual(os.DirFS(resultFolder), reader)
		assert.NoErrorf(t, err, "[example %d]", idx+1)
		assert.Truef(t, equal, "[example %d]", idx+1)
	}
}

func TestSpecInvalid(t *testing.T) {
	testFiles, err := getHrxFiles(invalidPath)
	require.NoError(t, err)

	t.Log(testFiles)
	assert.Positive(t, len(testFiles))

	for idx, fname := range testFiles {
		_, err := hrx.OpenReader(fname)

		assert.Errorf(t, err, "[example %d]", idx+1)
	}
}
