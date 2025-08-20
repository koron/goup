package install

import (
	"context"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestInstallZip(t *testing.T) {
	root := t.TempDir()
	got, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root, "-goos", "windows", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "", got)
	testutil.AssertGodir(t, root, "go1.18.6.windows-amd64")
}

func TestInstallTarGz(t *testing.T) {
	root := t.TempDir()
	got, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root, "-goos", "linux", "-goarch", "amd64", "go1.18.6"})
		if err != nil {
			t.Errorf("install failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "", got)
	testutil.AssertGodir(t, root, "go1.18.6.linux-amd64")
}
