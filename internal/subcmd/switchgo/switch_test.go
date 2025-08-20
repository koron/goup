package switchgo

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestSwitch(t *testing.T) {
	root := t.TempDir()
	godir0 := filepath.Join(root, "go1.18.6.windows-amd64")
	godir1 := filepath.Join(root, "go1.19.1.windows-amd64")
	testutil.AssertMkdirAll(t, godir0)
	testutil.AssertMkdirAll(t, godir1)
	testutil.AssertMkdirAll(t, filepath.Join(root, "go1.18.6.linux-amd64"))
	curr := filepath.Join(root, "current")
	// switch to go1.18.6
	out0, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "go1.18.6.windows-amd64\n", out0)
	testutil.AssertSymbolicLink(t, curr, "go1.18.6.windows-amd64")

	// switch to go1.19.1
	out1, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.19.1",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "go1.19.1.windows-amd64\n", out1)
	testutil.AssertSymbolicLink(t, curr, "go1.19.1.windows-amd64")
}

func TestSwitchUnmatch(t *testing.T) {
	root := t.TempDir()
	// switch to go1.18.6
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		testutil.AssertErr(t, err, `no installations for "go1.18.6"`)
	})
	testutil.AssertStdout(t, "", out)
}

func TestSwitchMatchMany(t *testing.T) {
	root := t.TempDir()
	// switch to go1.18.6
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6",
		})
		testutil.AssertErr(t, err, `no installations for "go1.18.6"`)
	})
	testutil.AssertStdout(t, "", out)
}

func TestSwitchDryrun(t *testing.T) {
	root := t.TempDir()
	godir := filepath.Join(root, "go1.18.6.windows-amd64")
	testutil.AssertMkdirAll(t, godir)
	curr := filepath.Join(root, "current")

	// switch to go1.18.6
	cout, cerr := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{
			"-root", root, "-goos", "windows", "-goarch", "amd64", "-dryrun", "go1.18.6",
		})
		if err != nil {
			t.Errorf("switch failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "go1.18.6.windows-amd64\n", cout)
	testutil.AssertStderr(t, "not installed because of dryrun\n", cerr)
	testutil.AssertIsNotExist(t, curr)
}
