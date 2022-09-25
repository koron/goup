package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func assertSymbolicLink(t *testing.T, name, dest string) {
	t.Helper()
	got, err := os.Readlink(name)
	if err != nil {
		t.Fatalf("failed to read a link: %s", err)
	}
	if want := dest; got != want {
		t.Fatalf("unexpected symbolic link: want=%s got=%s", want, got)
	}
}

func TestSwitch(t *testing.T) {
	root := t.TempDir()
	godir0 := filepath.Join(root, "go1.18.6.windows-amd64")
	godir1 := filepath.Join(root, "go1.19.1.windows-amd64")
	assertMkdirAll(t, godir0)
	assertMkdirAll(t, godir1)
	assertMkdirAll(t, filepath.Join(root, "go1.18.6.linux-amd64"))
	curr := filepath.Join(root, "current")
	// switch to go1.18.6
	out0, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("switch", flag.ContinueOnError)
		err := localSwitch(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	assertStdout(t, "go1.18.6.windows-amd64\n", out0)
	assertSymbolicLink(t, curr, "go1.18.6.windows-amd64")

	// switch to go1.19.1
	out1, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("switch", flag.ContinueOnError)
		err := localSwitch(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.19.1",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	assertStdout(t, "go1.19.1.windows-amd64\n", out1)
	assertSymbolicLink(t, curr, "go1.19.1.windows-amd64")
}

func TestSwitchUnmatch(t *testing.T) {
	root := t.TempDir()
	// switch to go1.18.6
	out, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("switch", flag.ContinueOnError)
		err := localSwitch(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		assertErr(t, err, `no installations for "go1.18.6"`)
	})
	assertStdout(t, "", out)
}

func TestSwitchMatchMany(t *testing.T) {
	root := t.TempDir()
	// switch to go1.18.6
	out, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("switch", flag.ContinueOnError)
		err := localSwitch(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		assertErr(t, err, `no installations for "go1.18.6"`)
	})
	assertStdout(t, "", out)
}

func TestSwitchDryrun(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	assertMkdirAll(t, godir)
	curr := filepath.Join(root, "current")

	// switch to go1.18.6
	cout, cerr := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("switch", flag.ContinueOnError)
		err := localSwitch(fs, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "-dryrun", "go1.18.6",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	assertStdout(t, "go1.18.6.windows-amd64\n", cout)
	assertStderr(t, "not installed because of dryrun\n", cerr)
	assertIsNotExist(t, curr)
}
