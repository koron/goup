package main

import (
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/goup/godlremote"
	"github.com/koron/goup/internal/dltestsrv"
)

func testSubcmd(t *testing.T, s *dltestsrv.Server, fn func()) string {
	return captureStdout(t, func() {
		withDltestsrv(t, s, fn)
	})
}

func withDltestsrv(t *testing.T, s *dltestsrv.Server, fn func()) {
	downloadBase := godlremote.DownloadBase
	if s == nil {
		s = &dltestsrv.Server{}
	}
	srv := httptest.NewServer(s)
	godlremote.DownloadBase = srv.URL
	defer func() {
		godlremote.DownloadBase = downloadBase
		srv.Close()
	}()
	fn()
}

func captureStdout(t *testing.T, fn func()) (out string) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Helper()
		t.Fatalf("failed take over os.Stdout: %s", err)
	}
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf strings.Builder
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			t.Helper()
			t.Errorf("goup testing: copying pipe: %s", err)
			return
		}
		outC <- buf.String()
	}()
	defer func() {
		w.Close()
		os.Stdout = stdout
		out = <-outC
	}()
	fn()
	return
}

func assertStdout(t *testing.T, want, got string) {
	d := cmp.Diff(want, got)
	if d != "" {
		t.Helper()
		t.Errorf("unexpected stdout: -want +got\n%s", d)
	}
}

func captureStderr(t *testing.T, fn func()) (out string) {
	stdout := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Helper()
		t.Fatalf("failed take over os.Stderr: %s", err)
	}
	os.Stderr = w
	outC := make(chan string)
	go func() {
		var buf strings.Builder
		_, err := io.Copy(&buf, r)
		r.Close()
		if err != nil {
			t.Helper()
			t.Errorf("goup testing: copying pipe: %s", err)
			return
		}
		outC <- buf.String()
	}()
	defer func() {
		w.Close()
		os.Stderr = stdout
		out = <-outC
	}()
	fn()
	return
}

func assertStderr(t *testing.T, want, got string) {
	d := cmp.Diff(want, got)
	if d != "" {
		t.Helper()
		t.Errorf("unexpected stderr: -want +got\n%s", d)
	}
}
