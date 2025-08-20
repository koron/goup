package uninstall

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/common"
)

var Command = subcmd.DefineCommand("uninstall", "uninstall Go releases", func(ctx context.Context, args []string) error {
	var (
		root   string
		goos   string
		goarch string
		clean  bool
	)
	fs := subcmd.FlagSet(ctx)
	fs.StringVar(&root, "root", common.GoupRoot(), "root dir to install")
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

	uni := Uninstaller{
		RootDir: root,
		GOOS:    goos,
		GOARCH:  goarch,
		Clean:   clean,
	}
	for _, ver := range versions {
		err := uni.Uninstall(ctx, ver)
		if err != nil {
			return err
		}
	}
	return nil
})

type Uninstaller struct {
	RootDir string
	GOOS    string
	GOARCH  string
	Clean   bool
}

func (uni Uninstaller) Uninstall(ctx context.Context, ver string) error {
	var deleted bool
	name := fmt.Sprintf("%s.%s-%s", ver, uni.GOOS, uni.GOARCH)
	dir := filepath.Join(uni.RootDir, name)
	ok, err := uni.hasDir(dir)
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
	if uni.Clean {
		pat := filepath.Join(uni.RootDir, "dl", name+".*")
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

func (uni Uninstaller) hasDir(name string) (bool, error) {
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
