package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgradeCmd(t *testing.T) {
	root := t.TempDir()
	fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	err := upgradeCmd(fs, []string{"-root", root})
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: check result
}

func TestUpgradeCmdEmptyRoot(t *testing.T) {
	_ = captureStderr(t, func() {
		fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
		err := upgradeCmd(fs, []string{"-root", ""})
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
	// an upgrade detected, it is "current"
	root := t.TempDir()
	goname := "go1.19.windows-amd64"
	err := os.MkdirAll(filepath.Join(root, goname), 0777)
	if err != nil {
		t.Errorf("mkdir failed: %v", err)
		return
	}
	err = switchGo(root, "current", goname)
	if err != nil {
		t.Errorf("switch failed: %v", err)
		return
	}

	got := captureStderr(t, func() {
		fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
		err = upgradeCmd(fs, []string{"-root", root, "-dryrun"})
		if err != nil {
			t.Error(err)
		}
	})
	assertStderr(t, strings.Join([]string{
		"upgraded Go go1.19.windows-amd64 to go1.19.1.windows-amd64",
		""}, "\n"), got)

	// FIXME: check result
}
