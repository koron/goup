package list

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestLocalList(t *testing.T) {
	root := t.TempDir()
	godir0 := filepath.Join(root, "go1.18.6.windows-amd64")
	godir1 := filepath.Join(root, "go1.19.1.windows-amd64")
	testutil.AssertMkdirAll(t, godir0)
	testutil.AssertMkdirAll(t, godir1)
	testutil.AssertSymlink(t, "go1.18.6.windows-amd64", filepath.Join(root, "current"))
	out, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root})
		if err != nil {
			t.Errorf("list failed: %s", err)
		}
	})
	testutil.AssertStdout(t, strings.Join([]string{
		"Local Version:",
		"  go1.19.1.windows-amd64",
		"  go1.18.6.windows-amd64 (current)",
		"",
	}, "\n"), out)
}
