package hrx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/liamg/memoryfs"
)

var boundaryRe = regexp.MustCompile("^<=+>")
var pathRe = regexp.MustCompile(`^[^/]+(/[^/]+)*/?$`)

type hrxArchive struct {
	fs *memoryfs.FS
}

// explicitly proxy the methods to make sure we don't leak any extra methods
// from memoryfs

func (h *hrxArchive) Open(name string) (fs.File, error) {
	return h.fs.Open(name)
}

func (h *hrxArchive) Glob(pattern string) ([]string, error) {
	return h.fs.Glob(pattern)
}

func (h *hrxArchive) ReadDir(name string) ([]os.DirEntry, error) {
	return h.fs.ReadDir(name)
}

func (h *hrxArchive) ReadFile(name string) ([]byte, error) {
	return h.fs.ReadFile(name)
}

func (h *hrxArchive) Stat(name string) (os.FileInfo, error) {
	return h.fs.Stat(name)
}

func (h *hrxArchive) Sub(dir string) (fs.FS, error) {
	return h.fs.Sub(dir)
}

// Open parses the target hrx file and either returns
// an fs compatible reader or an error
func OpenReader(fname string) (*hrxArchive, error) {
	memfs := memoryfs.New()

	fh, err := os.Open(fname)

	if err != nil {
		return nil, err
	}

	if err := ingestHrx(memfs, fh); err != nil {
		return nil, err
	}

	return &hrxArchive{
		fs: memfs,
	}, nil
}

type lineReader struct {
	reader *bufio.Reader
	buf    *string
	lineNr int
}

func newLineReader(r io.Reader) *lineReader {
	return &lineReader{
		reader: bufio.NewReader(r),
	}
}

func (lr *lineReader) PeekLine() (string, error) {
	if lr.buf != nil {
		return *lr.buf, nil
	}

	// error does not matter, we will get no data in case if EOF
	l, _ := lr.reader.ReadString('\n')

	if l != "" {
		lr.buf = &l
		return l, nil
	}

	return "", io.EOF
}

func (lr *lineReader) ReadLine() (string, error) {
	lr.lineNr++

	if lr.buf != nil {
		l := lr.buf
		lr.buf = nil
		return *l, nil
	}

	// error does not matter, we will get no data in case if EOF
	l, _ := lr.reader.ReadString('\n')

	if l != "" {
		return l, nil
	}

	return "", io.EOF
}

func ingestHrx(fs *memoryfs.FS, reader io.Reader) error {
	lr := newLineReader(reader)

	firstLine, err := lr.PeekLine()

	if err != nil {
		return fmt.Errorf("Failed to get the very first line")
	}

	boundary := boundaryRe.FindString(firstLine)

	if boundary == "" {
		return fmt.Errorf("Boundary has not been found on the first line")
	}

	for {
		err := ingestEntry(fs, lr, boundary)

		if err == io.EOF { // nothing else to parse
			return nil
		} else if err != nil {
			return err
		}
	}
}

func ingestEntry(fs *memoryfs.FS, lr *lineReader, boundary string) error {
	var current string
	var err error

	reportError := func(err2 error) error {
		return fmt.Errorf("[line %d] %w , line: [%s]", lr.lineNr, err2, strings.TrimRight(current, "\n"))
	}

	current, err = lr.ReadLine()

	if err != nil {
		return err
	}

	if !strings.HasPrefix(current, boundary) {
		return reportError(fmt.Errorf("new entry does not begin with a boundary, boundary: [%s]", boundary))
	}

	// slurp optional comment
	if strings.TrimRight(current, "\n") == boundary {
		for {
			current, err = lr.ReadLine()

			if err != nil {
				return err
			}

			if strings.HasPrefix(current, boundary) {
				break
			}
		}
	}

	parts := strings.SplitN(strings.TrimRight(current, "\n"), " ", 2)

	if len(parts) != 2 {
		return reportError(fmt.Errorf("entry does not have a path in it"))
	}

	p := parts[1]

	if !pathValid(p) {
		return reportError(fmt.Errorf("path is invalid [%s]", p))
	}

	if strings.HasSuffix(p, "/") {
		err := fs.MkdirAll(p, 0o700)

		if err != nil {
			return err
		}

		for {
			next, err := lr.PeekLine()

			if err != nil {
				break
			}

			if next != "\n" {
				break
			}

			current, err = lr.ReadLine()

			if err != nil {
				return err
			}
		}

		return nil
	}

	err = fs.MkdirAll(path.Dir(p), 0o700)

	if err != nil {
		return err
	}

	var b bytes.Buffer

	for {
		next, err := lr.PeekLine()

		if err != nil {
			break
		}

		if strings.HasPrefix(next, boundary) {
			break
		}

		curr, err := lr.ReadLine()

		if err != nil {
			return err
		}

		next, _ = lr.PeekLine()

		if strings.HasPrefix(next, boundary) {
			b.WriteString(strings.TrimRight(curr, "\n"))
		} else {
			b.WriteString(curr)
		}
	}

	return fs.WriteFile(p, b.Bytes(), 0o644)
}

func pathValid(p string) bool {
	if !pathRe.MatchString(p) {
		return false
	}

	for _, s := range strings.Split(p, "/") {
		if s == "." || s == ".." {
			return false
		}
	}

	return true
}
