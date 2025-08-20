package common

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/koron/goup/internal/verutil"
	"golang.org/x/mod/semver"
)

// GoupRoot retrieves GOUP_ROOT environment variable's value
func GoupRoot() string {
	s := os.Getenv("GOUP_ROOT")
	if s != "" {
		return s
	}
	// for comaptibility. this will be removed in future version.
	return os.Getenv("GODL_ROOT")
}

// GoupLinkname retrieves GOUP_LINKNAME environment variable's value
func GoupLinkname() string {
	s := os.Getenv("GOUP_LINKNAME")
	if s != "" {
		return s
	}
	return "current"
}

// LocalCurrent gets name of Go directory which is selected as "current"
// version.
func LocalCurrent(name string) (string, error) {
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

type InstalledGo struct {
	Version string
	OS      string
	Arch    string
	Name    string
	Semver  string
}

type InstalledGos []InstalledGo

func (list InstalledGos) Filter(f func(InstalledGo) bool) InstalledGos {
	var res InstalledGos
	for _, g := range list {
		if f(g) {
			res = append(res, g)
		}
	}
	return res
}

var rxGoDir = regexp.MustCompile(`^(go\d+(?:\.\d+)*(?:(?:rc|beta|alpha)\d+)?)\.(\D[^-]*)-(.+)$`)

// ListInstalledGo lists installed Go verions.
func ListInstalledGo(root string) (InstalledGos, error) {
	filist, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	list := make(InstalledGos, 0, len(filist))
	for _, fi := range filist {
		if !fi.IsDir() {
			continue
		}
		m := rxGoDir.FindStringSubmatch(fi.Name())
		if m == nil {
			continue
		}
		ver := verutil.Regulate(m[1])
		if ver == "" {
			continue
		}
		list = append(list, InstalledGo{
			Version: m[1],
			OS:      m[2],
			Arch:    m[3],
			Name:    m[0],
			Semver:  ver,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		cmp := semver.Compare(list[i].Semver, list[j].Semver)
		if cmp != 0 {
			return cmp > 0
		}
		return list[i].Semver > list[j].Semver
	})
	return list, nil
}

// SwitchGo switches/updates "current" symbolic link to goName.
func SwitchGo(root, linkname, goName string) error {
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
