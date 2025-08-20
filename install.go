package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/koron-go/zipx"
	"github.com/koron/goup/godlremote"
	"github.com/koron/goup/tarx"
)

func installCmd(ctx context.Context, args []string) error {
	var root string
	var force bool
	var goos string
	var goarch string
	fs := context2flagset(ctx)
	fs.StringVar(&root, "root", envGoupRoot(), "root dir to install")
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

// dldir creates/assures a directory to download archives.
func (ins installer) dldir() (string, error) {
	dir := filepath.Join(ins.rootdir, "dl")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// exdir creates/assures a directory to extracting an archive.
func (ins installer) extdir(f godlremote.File) (string, error) {
	dir := filepath.Join(ins.rootdir, f.Name())
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
	err = af.Download(ctx, name, ins.force)
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

func (ins installer) extract(ctx context.Context, dstdir string, srcfile string) error {
	if strings.HasSuffix(srcfile, ".zip") {
		err := zipx.New().ExtractFile(ctx, srcfile, zipDestDir(dstdir))
		if err != nil {
			return err
		}
		return nil
	}
	if strings.HasSuffix(srcfile, ".tar.gz") {
		err := tarx.ExtractFile(ctx, srcfile, tarDestDir(dstdir))
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unsupported archive: %s", srcfile)
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
	return &outFile{file: f, mtime: info.Modified}, nil
}

type tarDestDir string

func (d tarDestDir) CreateDir(info tarx.DirInfo) error {
	if !strings.HasPrefix(info.Name, "go/") {
		return fmt.Errorf("dir not under go/: %s", info.Name)
	}
	dir := filepath.Join(string(d), info.Name[3:])
	err := os.MkdirAll(dir, os.FileMode(info.Mode))
	if err != nil {
		return err
	}
	return nil
}

func (d tarDestDir) CreateFile(info tarx.FileInfo) (io.Writer, error) {
	if !strings.HasPrefix(info.Name, "go/") {
		return nil, fmt.Errorf("file not under go/: %s", info.Name)
	}
	info.Name = filepath.Join(string(d), info.Name[3:])
	err := os.MkdirAll(filepath.Dir(info.Name), 0777)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(info.Name)
	if err != nil {
		return nil, err
	}
	return &outFile{
		file:  f,
		mtime: info.ModTime,
		mode:  os.FileMode(info.Mode),
	}, nil
}

type outFile struct {
	file  *os.File
	mtime time.Time
	mode  os.FileMode
}

func (of outFile) Write(b []byte) (int, error) {
	return of.file.Write(b)
}

func (of outFile) Close() error {
	name := of.file.Name()
	err := of.file.Close()
	if err != nil {
		return err
	}
	if !of.mtime.IsZero() {
		err := os.Chtimes(name, of.mtime, of.mtime)
		if err != nil {
			return err
		}
	}
	if of.mode != 0 {
		err := os.Chmod(name, of.mode)
		if err != nil {
			return err
		}
	}
	return nil
}
