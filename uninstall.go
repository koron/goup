package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func uninstall(fs *flag.FlagSet, args []string) error {
	var root string
	var goos string
	var goarch string
	var clean bool
	fs.StringVar(&root, "root", os.Getenv("GODL_ROOT"), "root dir to install")
	fs.StringVar(&goos, "goos", runtime.GOOS, "OS for go to install")
	fs.StringVar(&goarch, "goarch", runtime.GOARCH, "ARCH for go to install")
	fs.BoolVar(&clean, "clean", false, "clean distfiles")
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
	uni := uninstaller{
		rootdir: root,
		goos:    goos,
		goarch:  goarch,
		clean:   clean,
	}
	for _, ver := range versions {
		err := uni.uninstall(ctx, ver)
		if err != nil {
			return err
		}
	}
	return nil
}

type uninstaller struct {
	rootdir string
	goos    string
	goarch  string
	clean   bool
}

func (uni uninstaller) uninstall(ctx context.Context, ver string) error {
	var deleted bool
	name := fmt.Sprintf("%s.%s-%s", ver, uni.goos, uni.goarch)
	dir := filepath.Join(uni.rootdir, name)
	ok ,err := uni.hasDir(dir)
	if err != nil {
		return err
	}
	if ok {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
		deleted = true
	}
	// remove distributed files.
	if uni.clean {
		pat := filepath.Join(uni.rootdir, "dl", name+".*")
		paths, err := filepath.Glob(pat)
		if err != nil {
			return err
		}
		for _, p := range paths {
			err := os.Remove(p)
			if err != nil {
				return err
			}
			deleted = true
		}
	}
	if !deleted {
		return fmt.Errorf("no deleted files for %s", name)
	}
	return nil
}

func (uni uninstaller) hasDir(name string) (bool , error) {
	fi, err := os.Lstat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("not directory: %s", name)
	}
	return true, nil
}
