package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallZip(t *testing.T) {
	root := t.TempDir()
	got := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("install", flag.ContinueOnError)
		err := install(fs, []string{"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	assertStdout(t, "", got)

	// check installed Go directory
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	fi, err := os.Stat(godir)
	if err != nil {
		t.Errorf("failed to stat: %s", err)
		return
	}
	if !fi.IsDir() {
		t.Errorf("not found install directory: %s", godir)
	}
}

func TestInstallTarGz(t *testing.T) {
	root := t.TempDir()
	got := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("install", flag.ContinueOnError)
		err := install(fs, []string{"-root", root, "-goos", "linux", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	assertStdout(t, "", got)

	// check installed Go directory
	godir := filepath.Join(root, "go1.18.6.linux-amd64")
	fi, err := os.Stat(godir)
	if err != nil {
		t.Errorf("failed to stat: %s", err)
		return
	}
	if !fi.IsDir() {
		t.Errorf("not found install directory: %s", godir)
	}
}
