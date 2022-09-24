package main

import (
	"flag"
	"testing"
)

func TestInstallZip(t *testing.T) {
	root := t.TempDir()
	got := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("install", flag.ContinueOnError)
		err := installCmd(fs, []string{"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	assertStdout(t, "", got)
	assertGodir(t, root, "go1.18.6.windows-amd64")
}

func TestInstallTarGz(t *testing.T) {
	root := t.TempDir()
	got := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("install", flag.ContinueOnError)
		err := installCmd(fs, []string{"-root", root, "-goos", "linux", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	assertStdout(t, "", got)
	assertGodir(t, root, "go1.18.6.linux-amd64")
}
