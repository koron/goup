package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/koron-go/zipx"
	"github.com/koron/godltool/godlremote"
)

func install(fs *flag.FlagSet, args []string) error {
	var root string
	var force bool
	var goos string
	var goarch string
	fs.StringVar(&root, "root", os.Getenv("GODL_ROOT"), "root dir to install")
	fs.BoolVar(&force, "force", false, "override installation")
	fs.StringVar(&goos, "goos", runtime.GOOS, "OS for go to install")
	fs.StringVar(&goarch, "goarch", runtime.GOARCH, "ARCH for go to install")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}
	versions := fs.Args()
	if len(versions) == 0 {
		fs.PrintDefaults()
		return errors.New("no versions to install")
	}

	ctx := context.Background()
	rels, err := godlremote.Download(ctx, true)
	if err != nil {
		return err
	}
	ins := installer{
		releases: rels,
		rootdir:  root,
		force:    force,
		goos:     goos,
		goarch:   goarch,
	}
	for _, ver := range versions {
		err := ins.install(ctx, ver)
		if err != nil {
			return err
		}
	}
	return nil
}

type installer struct {
	releases godlremote.Releases
	rootdir  string
	force    bool
	goos     string
	goarch   string
}

func (ins installer) archiveFile(ver string) (godlremote.File, bool) {
	for _, r := range ins.releases {
		if r.Version != ver {
			continue
		}
		for _, f := range r.Files {
			if f.Kind != "archive" || f.OS != ins.goos || f.Arch != ins.goarch {
				continue
			}
			return f, true
		}
	}
	return godlremote.File{}, false
}

func (ins installer) dldir() (string, error) {
	dir := filepath.Join(ins.rootdir, "dl")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (ins installer) extdir(f godlremote.File) (string, error) {
	dir := filepath.Join(ins.rootdir, fmt.Sprintf("%s.%s-%s", f.Version, f.OS, f.Arch))
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (ins installer) install(ctx context.Context, ver string) error {
	af, ok := ins.archiveFile(ver)
	if !ok {
		return fmt.Errorf("no archives found for version=%s OS=%s arch=%s",
			ver, ins.goos, ins.goarch)
	}
	dldir, err := ins.dldir()
	if err != nil {
		return err
	}
	name := filepath.Join(dldir, af.Filename)
	err = af.Download(ctx, name)
	if err != nil {
		return err
	}
	extdir, err := ins.extdir(af)
	if err != nil {
		return err
	}
	err = ins.extract(ctx, extdir, name)
	if err != nil {
		os.RemoveAll(extdir)
		return err
	}
	return nil
}

func (ins installer) extract(ctx context.Context, distdir string, srcfile string) error {
	err := zipx.New().ExtractFile(ctx, srcfile, zipDestDir(distdir))
	if err != nil {
		return err
	}
	return nil
}

type zipDestDir string

func (d zipDestDir) CreateDir(name string, info zipx.DirInfo) error {
	if !strings.HasPrefix(name, "go/") {
		return fmt.Errorf("dir not under go/: %s", name)
	}
	dir := filepath.Join(string(d), name[3:])
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (d zipDestDir) CreateFile(name string, info zipx.FileInfo) (io.Writer, error) {
	if !strings.HasPrefix(name, "go/") {
		return nil, fmt.Errorf("file not under go/: %s", name)
	}
	name = filepath.Join(string(d), name[3:])
	err := os.MkdirAll(filepath.Dir(name), 0777)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &zipFile{File: f, mtime: info.Modified}, nil
}

type zipFile struct {
	*os.File
	mtime time.Time
}

func (zf zipFile) Write(b []byte) (int, error) {
	return zf.File.Write(b)
}

func (zf zipFile) Close() error {
	err := zf.File.Close()
	if err != nil {
		return err
	}
	return os.Chtimes(zf.File.Name(), zf.mtime, zf.mtime)
}
