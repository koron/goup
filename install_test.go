package main

import (
	"flag"
	"testing"
)

func TestInstall(t *testing.T) {
	got := testSubcmd(t, nil, func() {
		root := t.TempDir()
		fs := flag.NewFlagSet("install", flag.ContinueOnError)
		err := install(fs, []string{"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	assertStdout(t, "", got)
}
