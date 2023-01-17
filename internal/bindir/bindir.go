package bindir

import (
	"errors"
	"os"
	"path/filepath"
)

type Bindir struct {
	name string
	dir  *os.File
}

func gobin() (string, error) {
	gobin := os.Getenv("GOBIN")
	if gobin != "" {
		return gobin, nil
	}
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return filepath.Join(gopath, "bin"), nil
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, "go", "bin"), nil
}

func Open() (*Bindir, error) {
	name, err := gobin()
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(name) {
		return nil, errors.New("cannot operate, GOBIN must be an absolute path")
	}
	dir, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Bindir{
		name: name,
		dir:  dir,
	}, nil
}

func (b *Bindir) Close() {
	if b.dir != nil {
		b.dir.Close()
		b.dir = nil
	}
}

func (b *Bindir) Read() (string, error) {
	for {
		ents, err := b.dir.ReadDir(1)
		if err != nil {
			return "", err
		}
		ent := ents[0]
		if ent.IsDir() {
			continue
		}
		mode := ent.Type()
		if mode.IsRegular(){
			return filepath.Join(b.name, ent.Name()), nil
		}
	}
}
