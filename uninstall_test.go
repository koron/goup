package main

import (
	"flag"
	"path/filepath"
	"testing"
)

func TestUninstall(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	assertMkdirAll(t, godir)
	// uninstall
	out := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
		err := uninstallCmd(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		if err != nil {
			t.Errorf("uninstall failed: %s", err)
		}
	})
	assertStdout(t, "", out)
	assertIsNotExist(t, godir)
}

func TestUninstallInvalid(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	assertIsNotExist(t, godir)
	// uninstall
	out := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
		err := uninstallCmd(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		assertErr(t, err, "no deleted files for go1.18.6.windows-amd64")
	})
	assertStdout(t, "", out)
	assertIsNotExist(t, godir)
}

func TestUninstallClean(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	assertMkdirAll(t, godir)
	dlfile := filepath.Join(root, "dl", "go1.18.6.windows-amd64.zip")
	assertTouchFile(t, dlfile)
	// uninstall
	out := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
		err := uninstallCmd(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "-clean", "go1.18.6",
		})
		if err != nil {
			t.Errorf("uninstall failed: %s", err)
		}
	})
	assertStdout(t, "", out)
	assertIsNotExist(t, godir)
	assertIsNotExist(t, dlfile)
}
