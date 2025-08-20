package testutil

import (
	"context"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/goup/godlremote"
	"github.com/koron/goup/internal/dltestsrv"
)

func TestSubcmd(t *testing.T, s *dltestsrv.Server, fn func(context.Context)) (capturedOut, capturedErr string) {
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

func AssertStdout(t *testing.T, want, got string) {
	d := cmp.Diff(want, got)
	if d != "" {
		t.Helper()
		t.Errorf("unexpected stdout: -want +got\n%s", d)
	}
}
