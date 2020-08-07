package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var rxGoArchive = regexp.MustCompile(`^(go\d+(?:\.\d+)*(?:(?:rc|beta|alpha)\d+)?\.\D[^-]*-.+)\.(?:zip|tar\.gz)$`)

func localClean(fs *flag.FlagSet, args []string) error {
	var root string
	var dryrun bool
	var all bool
	fs.StringVar(&root, "root", os.Getenv("GODL_ROOT"), "root dir to install")
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
	filist, err := ioutil.ReadDir(dldir)
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
}
