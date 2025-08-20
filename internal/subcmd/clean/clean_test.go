package clean

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestLocalClean(t *testing.T) {
	root := t.TempDir()
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"))
	testutil.AssertMkdirAll(t, filepath.Join(root, "go1.19.6.windows-amd64"))
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root})
		if err != nil {
			t.Errorf("clean failed: %s", err)
		}
	})
	testutil.AssertStdout(t, "", out)
	testutil.AssertIsNotExist(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	testutil.AssertIsNotExist(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	testutil.AssertIsExist(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"), false)
}

func TestLocalCleanDryrun(t *testing.T) {
	root := t.TempDir()
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"))
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"))
	testutil.AssertTouchFile(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"))
	testutil.AssertMkdirAll(t, filepath.Join(root, "go1.19.6.windows-amd64"))
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root, "-dryrun"})
		if err != nil {
			t.Errorf("clean failed: %s", err)
		}
	})
	testutil.AssertStdout(t, strings.Join([]string{
		filepath.Join(root, "dl", "go1.18.windows-amd64.zip") + " will be deleted",
		filepath.Join(root, "dl", "go1.19.windows-amd64.zip") + " will be deleted",
		"",
	}, "\n"), out)
	testutil.AssertIsExist(t, filepath.Join(root, "dl", "go1.18.windows-amd64.zip"), false)
	testutil.AssertIsExist(t, filepath.Join(root, "dl", "go1.19.windows-amd64.zip"), false)
	testutil.AssertIsExist(t, filepath.Join(root, "dl", "go1.19.6.windows-amd64.zip"), false)
}
