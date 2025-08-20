package uninstall

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestUninstall(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	testutil.AssertMkdirAll(t, godir)
	// uninstall
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		if err != nil {
			t.Errorf("uninstall failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "", out)
	testutil.AssertIsNotExist(t, godir)
}

func TestUninstallInvalid(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	testutil.AssertIsNotExist(t, godir)
	// uninstall
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		testutil.AssertErr(t, err, "no deleted files for go1.18.6.windows-amd64")
	})
	testutil.AssertStdout(t, "", out)
	testutil.AssertIsNotExist(t, godir)
}

func TestUninstallClean(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	testutil.AssertMkdirAll(t, godir)
	dlfile := filepath.Join(root, "dl", "go1.18.6.windows-amd64.zip")
	testutil.AssertTouchFile(t, dlfile)
	// uninstall
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "-clean", "go1.18.6",
		})
		if err != nil {
			t.Errorf("uninstall failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "", out)
	testutil.AssertIsNotExist(t, godir)
	testutil.AssertIsNotExist(t, dlfile)
}
