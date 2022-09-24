package main

import (
	"flag"
	"strings"
	"testing"
)

func TestRemoteList(t *testing.T) {
	got, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("remotelist", flag.ContinueOnError)
		err := remoteList(fs, nil)
		if err != nil {
			t.Errorf("remoteList failed: %s", err)
		}
	})
	assertStdout(t, strings.Join([]string{
		"Remote Version:",
		"  go1.19.1",
		"  go1.18.6",
		""}, "\n"), got)
}

func TestRemoteListMatch(t *testing.T) {
	got, _ := testSubcmd(t, nil, func() {
		fs := flag.NewFlagSet("remotelist", flag.ContinueOnError)
		err := remoteList(fs, []string{"-all", "-match", "1\\.18"})
		if err != nil {
			t.Errorf("remoteList failed: %s", err)
		}
	})
	assertStdout(t, strings.Join([]string{
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
