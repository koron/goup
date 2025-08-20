package main

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalClean(t *testing.T) {
	root := t.TempDir()
	assertTouchFile(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	assertTouchFile(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	assertTouchFile(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"))
	assertMkdirAll(t, filepath.Join(root, "go1.19.6.windows-amd64"))
	out, _ := testSubcmd(t, nil, func(ctx context.Context) {
		err := cleanCommand.Run(ctx, []string{"-root", root})
		if err != nil {
			t.Errorf("clean failed: %s", err)
		}
	})
	assertStdout(t, "", out)
	assertIsNotExist(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	assertIsNotExist(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	assertIsExist(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"), false)
}

func TestLocalCleanDryrun(t *testing.T) {
	root := t.TempDir()
	assertTouchFile(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	assertTouchFile(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	assertTouchFile(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"))
	assertMkdirAll(t, filepath.Join(root, "go1.19.6.windows-amd64"))
	out, _ := testSubcmd(t, nil, func(ctx context.Context) {
		err := cleanCommand.Run(ctx, []string{"-root", root, "-dryrun"})
		if err != nil {
			t.Errorf("clean failed: %s", err)
		}
	})
	assertStdout(t, strings.Join([]string{
		filepath.Join(root, "dl", "go1.18.windows-amd64.zip") + " will be deleted",
		filepath.Join(root, "dl", "go1.19.windows-amd64.zip") + " will be deleted",
		"",
	}, "\n"), out)
	assertIsExist(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"), false)
	assertIsExist(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"), false)
	assertIsExist(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"), false)
}
