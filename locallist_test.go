package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assertSymlink(t *testing.T, oldname, newname string) {
	err := os.Symlink(oldname, newname)
	if err != nil {
		t.Helper()
		t.Fatalf("failed to create symbolic link: %s", err)
	}
}

func TestLocalList(t *testing.T) {
	root := t.TempDir()
	godir0 := filepath.Join(root, "go1.18.6.windows-amd64")
	godir1 := filepath.Join(root, "go1.19.1.windows-amd64")
	assertMkdirAll(t, godir0)
	assertMkdirAll(t, godir1)
	assertSymlink(t, "go1.18.6.windows-amd64", filepath.Join(root, "current"))
	out, _ := testSubcmd(t, nil, func(ctx context.Context) {
		err := listCommand.Run(ctx, []string{"-root", root})
		if err != nil {
			t.Errorf("list failed: %s", err)
		}
	})
	assertStdout(t, strings.Join([]string{
		"Local Version:",
		"  go1.19.1.windows-amd64",
		"  go1.18.6.windows-amd64 (current)",
		"",
	}, "\n"), out)
}
