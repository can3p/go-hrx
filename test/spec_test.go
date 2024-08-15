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

	onlyFile := os.Getenv("ONLY_FILE")
	onlyFileFound := false

	for _, fname := range testFiles {
		if onlyFile != "" {
			if onlyFile != fname {
				continue
			}
			onlyFileFound = true
		}

		dirname, found := strings.CutSuffix(fname, ".hrx")
		require.Truef(t, found, "[example %s]", fname)

		fullHrxPath := path.Join(hrxPath, fname)
		resultFolder := path.Join(extractedPath, dirname)

		reader, err := hrx.OpenReader(fullHrxPath)

		assert.NoErrorf(t, err, "[example %s]", fullHrxPath)

		if err != nil {
			continue
		}

		err = DirsEqual(os.DirFS(resultFolder), reader)
		assert.NoErrorf(t, err, "[example %s]", fullHrxPath)
	}

	require.True(t, onlyFile == "" || onlyFileFound, "ONLY_FILE environment is set, but no file matches it")
}

func TestSpecInvalid(t *testing.T) {
	testFiles, err := getHrxFiles(invalidPath)
	require.NoError(t, err)

	t.Log(testFiles)
	assert.Positive(t, len(testFiles))

	for _, fname := range testFiles {
		fullHrxPath := path.Join(hrxPath, fname)

		_, err := hrx.OpenReader(fullHrxPath)

		assert.Errorf(t, err, "[example %s]", fullHrxPath)
	}
}
