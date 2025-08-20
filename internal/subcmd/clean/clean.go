// Package clean provides "clean" subcmd of goup.
package clean

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/koron-go/subcmd"
	"github.com/koron/goup/internal/common"
)

var rxGoArchive = regexp.MustCompile(`^(go\d+(?:\.\d+)*(?:(?:rc|beta|alpha)\d+)?\.\D[^-]*-.+)\.(?:zip|tar\.gz)$`)

var Command = subcmd.DefineCommand("clean", "clean download caches", func(ctx context.Context, args []string) error {
	var root string
	var dryrun bool
	var all bool
	fs := subcmd.FlagSet(ctx)
	fs.StringVar(&root, "root", common.GoupRoot(), "root dir to install")
	fs.BoolVar(&dryrun, "dryrun", false, "don't switch, just test")
	fs.BoolVar(&all, "all", false, "clean all caches")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}

	dldir := filepath.Join(root, "dl")
	filist, err := os.ReadDir(dldir)
	if err != nil {
		return err
	}
	for _, fi := range filist {
		m := rxGoArchive.FindStringSubmatch(fi.Name())
		if m == nil {
			continue
		}
		if !all {
			// check the install dir existing or not.
			_, err := os.Stat(filepath.Join(root, m[1]))
			if err == nil {
				// existing
				continue
			}
			if !os.IsNotExist(err) {
				// general I/O error
				return err
			}
			// not found
		}
		name := filepath.Join(dldir, fi.Name())
		if dryrun {
			fmt.Printf("%s will be deleted\n", name)
			continue
		}
		err := os.RemoveAll(name)
		if err != nil {
			return err
		}
	}
	return nil
})
