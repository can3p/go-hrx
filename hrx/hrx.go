package hrx

import (
	"fmt"
	"io/fs"
)

// Open parses the target hrx file and either returns
// an fs compatible reader or an error
// TODO: add implementation
func OpenReader(fname string) (fs.FS, error) {
	return nil, fmt.Errorf("not implemented")
}
