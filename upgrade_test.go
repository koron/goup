package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func goName(base string) string {
	return fmt.Sprintf("%s.%s-%s", base, runtime.GOOS, runtime.GOARCH)
}

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
	fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
	err := upgradeCmd(fs, []string{"-root", ""})
	if err == nil {
		t.Fatal("want error but got no errors")
	}
	if s := err.Error() ; s != "required -root" {
		t.Errorf("unexpected error got=%s", s)
	}
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
	goname := goName("go1.19")
	err := os.MkdirAll(filepath.Join(root, goname), 0777)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("mkdir failed: %v", err)
	}
	err = switchGo(root, "current", goname)
	if err != nil {
		t.Fatalf("switch failed: %v", err)
	}
	err = upgrade(context.Background(), root, "current", true, false)
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: check result
}
