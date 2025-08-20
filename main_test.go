package main

import (
	"context"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/goup/godlremote"
	"github.com/koron/goup/internal/dltestsrv"
)

func testSubcmd(t *testing.T, s *dltestsrv.Server, fn func(context.Context)) (capturedOut, capturedErr string) {
	t.Helper()
	return captureStdoutStderr(t, func() {
		t.Helper()
		if s == nil {
			s = &dltestsrv.Server{}
		}
		srv := httptest.NewServer(s)
		defer srv.Close()
		ctx := godlremote.WithDownloadBase(context.Background(), srv.URL)
		fn(ctx)
	})
}

func captureStdoutStderr(t *testing.T, fn func()) (capturedOut, capturedErr string) {
	t.Helper()

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed take over os.Stdout: %s", err)
	}
	stdout := os.Stdout
	os.Stdout = outW
	outC := make(chan string)
	go func() {
		var buf strings.Builder
		_, err := io.Copy(&buf, outR)
		outR.Close()
		if err != nil {
			t.Helper()
			t.Errorf("goup testing: copying STDOUT pipe: %s", err)
			return
		}
		outC <- buf.String()
	}()
	defer func() {
		outW.Close()
		os.Stdout = stdout
		capturedOut = <-outC
	}()

	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed take over os.Stderr: %s", err)
	}
	stderr := os.Stderr
	os.Stderr = errW
	errC := make(chan string)
	go func() {
		var buf strings.Builder
		_, err := io.Copy(&buf, errR)
		errR.Close()
		if err != nil {
			t.Helper()
			t.Errorf("goup testing: copying STDERR pipe: %s", err)
			return
		}
		errC <- buf.String()
	}()
	defer func() {
		errW.Close()
		os.Stderr = stderr
		capturedErr = <-errC
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

func assertStderr(t *testing.T, want, got string) {
	d := cmp.Diff(want, got)
	if d != "" {
		t.Helper()
		t.Errorf("unexpected stderr: -want +got\n%s", d)
	}
}

// assertGodir checkk installed Go directory
func assertGodir(t *testing.T, root, goname string) {
	t.Helper()
	godir := filepath.Join(root, goname)
	assertIsExist(t, godir, true)
	assertIsExist(t, filepath.Join(godir, "README.txt"), false)
	// more checks in future. need to files and dirs adding to dummy archives.
}

func assertErr(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatal("an operation is succeeded, unexpectedly")
	}
	if got := err.Error(); want != got {
		t.Fatalf("an operation is failed with unexpected error:\nwant=%s\ngot=%s", want, got)
	}
}

// assertIsNotExist checks a file/dir is not exist
func assertIsNotExist(t *testing.T, name string) {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return
	}
	t.Helper()
	if err == nil {
		t.Fatalf("a file/dir is exist, unexpectedly: %s", name)
	}
	t.Fatalf("unexpected os.Stat failure: %s", err)
}

// assertIsExist checks a file/dir is exist
func assertIsExist(t *testing.T, name string, isDir bool) {
	t.Helper()
	fi, err := os.Stat(name)
	if err != nil {
		t.Fatalf("unexpected os.Stat failure: %s", err)
	}
	if want, got := isDir, fi.IsDir(); want != got {
		if want {
			t.Fatalf("a path %s is not a directory", name)
		} else {
			t.Fatalf("a path %s is not a file", name)
		}
	}
}

// assertMkdirAll creates a dir with parent directories.
func assertMkdirAll(t *testing.T, name string) {
	t.Helper()
	err := os.MkdirAll(name, 0755)
	if err != nil {
		t.Fatalf("failed to make directory: %s", err)
	}
	assertIsExist(t, name, true)
}

// assertTouchFile creates a file with parent directories.
func assertTouchFile(t *testing.T, name string) {
	t.Helper()
	dir := filepath.Dir(name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("failed to create parent directories %s: %s", dir, err)
	}
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("failed to create a file %s: %s", name, err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("failed to close a file %s: %s", name, err)
	}
	assertIsExist(t, name, false)
}
