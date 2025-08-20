package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/goup/internal/common"
	"github.com/koron/goup/internal/testutil"
)

func TestUpgradeCmd(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()
	err := Command.Run(ctx, []string{"-root", root})
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: check result
}

func TestUpgradeCmdEmptyRoot(t *testing.T) {
	_, _ = testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", ""})
		if err == nil {
			t.Errorf("want error but got no errors")
			return
		}
		if s := err.Error(); s != "required -root" {
			t.Errorf("unexpected error got=%s", s)
		}
	})
	// don't check help output
}

func TestUpgradeDryrun0(t *testing.T) {
	// no upgrades detected
	root := t.TempDir()
	err := upgrade(context.Background(), root, "current", true, false)
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: check result
}

func TestUpgradeDryrun1(t *testing.T) {
	// should detect an upgrade for "current" version.

	const goname = "go1.18.windows-amd64"

	root := t.TempDir()
	err := os.MkdirAll(filepath.Join(root, goname), 0777)
	if err != nil {
		t.Errorf("mkdir failed: %v", err)
		return
	}
	err = common.SwitchGo(root, "current", goname)
	if err != nil {
		t.Errorf("switch failed: %v", err)
		return
	}

	_, got := testutil.TestSubcmd(t, nil, func(ctx context.Context) {
		err := Command.Run(ctx, []string{"-root", root, "-dryrun"})
		if err != nil {
			t.Errorf("upgrade failed: %s", err)
		}
	})
	testutil.AssertStderr(t, strings.Join([]string{
		"upgraded Go go1.18.windows-amd64 to go1.18.6.windows-amd64",
		""}, "\n"), got)

	// FIXME: check result
}

func TestDebugInstalledGos(t *testing.T) {
	list := common.InstalledGos{
		common.InstalledGo{Version: "1", OS: "windows", Arch: "amd64", Name: "foo"},
		common.InstalledGo{Version: "2", OS: "windows", Arch: "amd64", Name: "bar"},
		common.InstalledGo{Version: "3", OS: "windows", Arch: "amd64", Name: "baz"},
		common.InstalledGo{Version: "4", OS: "windows", Arch: "amd64", Name: "qux"},
	}
	want := strings.Join([]string{
		"",
		"\tfoo",
		"\tbar",
		"\tbaz",
		"\tqux",
	}, "\n")
	got := debugInstalledGos(list).String()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("unexpected: -want +got\n%s", d)
	}
}

func TestDebugLatestReleases(t *testing.T) {
	rels := map[string]latestRelease{
		"v1.18": {semver: "v1.18.6"},
		"v1.19": {semver: "v1.19.1"},
	}
	want := strings.Join([]string{
		"",
		"\tv1.19.1",
		"\tv1.18.6",
	}, "\n")
	got := debugLatestReleases(rels).String()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("unexpected: -want +got\n%s", d)
	}
}
