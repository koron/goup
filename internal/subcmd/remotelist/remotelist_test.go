package remotelist

import (
	"context"
	"strings"
	"testing"

	"github.com/koron/goup/internal/testutil"
)

func TestRemoteList(t *testing.T) {
	got, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, nil)
		if err != nil {
			t.Errorf("remoteList failed: %s", err)
		}
	})
	testutil.AssertStdout(t, strings.Join([]string{
		"Remote Version:",
		"  go1.19.1",
		"  go1.18.6",
		""}, "\n"), got)
}

func TestRemoteListMatch(t *testing.T) {
	got, _ := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-all", "-match", "1\\.18"})
		if err != nil {
			t.Errorf("remoteList failed: %s", err)
		}
	})
	testutil.AssertStdout(t, strings.Join([]string{
		"Remote Version:",
		"  go1.18.6",
		"  go1.18.5",
		"  go1.18.4",
		"  go1.18.3",
		"  go1.18.2",
		"  go1.18.1",
		"  go1.18",
		"  go1.18rc1",
		"  go1.18beta2",
		"  go1.18beta1",
		""}, "\n"), got)
}
