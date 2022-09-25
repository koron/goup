package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// localSwitch switches "current" selected Go version.
func localSwitch(fs *flag.FlagSet, args []string) error {
	var root string
	var goos string
	var goarch string
	var dryrun bool
	var linkname string
	fs.StringVar(&root, "root", envGoupRoot(), "root dir to install")
	fs.StringVar(&goos, "goos", runtime.GOOS, "OS for go to install")
	fs.StringVar(&goarch, "goarch", runtime.GOARCH, "ARCH for go to install")
	fs.BoolVar(&dryrun, "dryrun", false, "don't switch, just test")
	fs.StringVar(&linkname, "linkname", envGoupLinkname(), "name of symbolic link to switch")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}
	if n := fs.NArg(); n != 1 {
		fs.PrintDefaults()
		return fmt.Errorf("target to switch must be only one: %d", n)
	}
	target := fs.Arg(0)

	list, err := listInstalledGo(root)
	if err != nil {
		return err
	}
	list = list.filter(func(g installedGo) bool {
		if g.os != goos || g.arch != goarch {
			return false
		}
		return g.version == target
	})

	switch len(list) {
	case 0:
		return fmt.Errorf("no installations for %q", target)
	case 1:
		// nothing
	default:
		panic(fmt.Sprintf("hit %d installations for %q", len(list), target))
	}
	g := list[0]
	fmt.Println(g.name)
	if dryrun {
		fmt.Fprintln(os.Stderr, "not installed because of dryrun")
		return nil
	}
	return switchGo(root, linkname, g.name)
}

// switchGo switches/updates "current" symbolic link to goName.
func switchGo(root, linkname, goName string) error {
	dstdir := filepath.Join(root, linkname)
	// remove dstdir (symbolic link)
	_, err := os.Lstat(dstdir)
	if err == nil {
		err := os.Remove(dstdir)
		if err != nil {
			return err
		}
	}
	// create a symbolic link
	err = os.Symlink(goName, dstdir)
	if err != nil {
		return err
	}
	return nil
}
