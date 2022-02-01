package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"golang.org/x/mod/semver"
)

func localList(fs *flag.FlagSet, args []string) error {
	var root string
	var linkname string
	fs.StringVar(&root, "root", envGoupRoot(), "root dir to install")
	fs.StringVar(&linkname, "linkname", "current", "name of symbolic link to switch")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if root == "" {
		fs.PrintDefaults()
		return errors.New("required -root")
	}

	curr, err := localCurrent(filepath.Join(root, linkname))
	if err != nil {
		return err
	}

	list, err := listInstalledGo(root)
	if err != nil {
		return err
	}
	fmt.Println("Local Version:")
	for _, g := range list {
		if curr != "" && g.name == curr {
			fmt.Printf("  %s (%s)\n", g.name, linkname)
			continue
		}
		fmt.Printf("  %s\n", g.name)
	}
	return nil
}

type installedGo struct {
	version string
	os      string
	arch    string
	name    string
	semver  string
}

type installedGos []installedGo

func localCurrent(name string) (string, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return "", nil
	}
	rname, err := filepath.EvalSymlinks(name)
	if err != nil {
		return "", err
	}
	return filepath.Base(rname), nil
}

func (list installedGos) filter(f func(installedGo) bool) installedGos {
	var res installedGos
	for _, g := range list {
		if f(g) {
			res = append(res, g)
		}
	}
	return res
}

var rxGoDir = regexp.MustCompile(`^(go\d+(?:\.\d+)*(?:(?:rc|beta|alpha)\d+)?)\.(\D[^-]*)-(.+)$`)

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
		ver := regulateVersion(m[1])
		if ver == "" {
			continue
		}
		list = append(list, installedGo{
			version: m[1],
			os:      m[2],
			arch:    m[3],
			name:    m[0],
			semver:  ver,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		cmp := semver.Compare(list[i].semver, list[j].semver)
		if cmp != 0 {
			return cmp > 0
		}
		return list[i].semver > list[j].semver
	})
	return list, nil
}

var rxGoVer = regexp.MustCompile(`^go(?P<major>\d+)(?:\.(?P<minor>\d+)(?:\.(?P<patch>\d+))?)?(?P<pr>[A-Za-z][-.0-9A-Za-z]*)?`)

func regnum(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

func regulateVersion(s string) string {
	m := rxGoVer.FindStringSubmatch(s)
	if m == nil {
		return ""
	}
	bb := &bytes.Buffer{}
	bb.WriteRune('v')
	bb.WriteString(regnum(m[1]))
	bb.WriteRune('.')
	bb.WriteString(regnum(m[2]))
	bb.WriteRune('.')
	bb.WriteString(regnum(m[3]))
	if m[4] != "" {
		bb.WriteRune('-')
		bb.WriteString(m[4])
	}
	return bb.String()
}
