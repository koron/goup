package tarx

import (
	"io"
	"io/ioutil"
)

// Discard is a destination which discard all extracted files and dirs.
var Discard Destination = discard{}

type discard struct{}

func (discard) CreateDir(info DirInfo) error {
	return nil
}

func (discard) CreateFile(info FileInfo) (io.Writer, error) {
	return ioutil.Discard, nil
}
