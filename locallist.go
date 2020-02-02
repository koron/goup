package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func localList(fs *flag.FlagSet, args []string) error {
	var root string
	fs.StringVar(&root, "root", os.Getenv("GODL_ROOT"), "root dir to install")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}
	list, err := listInstalledGo(root)
	if err != nil {
		return err
	}
	fmt.Println("Local Version:")
	for _, g := range list {
		fmt.Printf("  %s\n", g.name)
	}
	return nil
}

type installedGo struct {
	version string
	os      string
	arch    string
	name    string
}

type installedGos []installedGo

func (list installedGos) filter(f func(installedGo) bool) installedGos {
	var res installedGos
	for _, g := range list {
		if f(g) {
			res = append(res, g)
		}
	}
	return res
}

var rxGoDir = regexp.MustCompile(`^(go\d+(?:\.\d+)*)\.(\D[^-]*)-(.+)$`)

func listInstalledGo(root string) (installedGos, error) {
	filist, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	list := make(installedGos, 0, len(filist))
	for _, fi := range filist {
		if !fi.IsDir() {
			continue
		}
		m := rxGoDir.FindStringSubmatch(fi.Name())
		if m == nil {
			continue
		}
		list = append(list, installedGo{
			version: m[1],
			os:      m[2],
			arch:    m[3],
			name:    m[0],
		})
	}
	// TODO: sort
	return list, nil
}
