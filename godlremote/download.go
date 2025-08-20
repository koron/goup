package godlremote

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const downloadBase = "https://go.dev/dl/"

type downloadBaseKey struct{}

func getBase(ctx context.Context) string {
	if overriddenBase, ok := ctx.Value(downloadBaseKey{}).(string); ok {
		return overriddenBase
	}
	return downloadBase
}

// WithDownloadBase overrides the base URL for downloads. Mainly for testing.
func WithDownloadBase(ctx context.Context, newBase string) context.Context {
	return context.WithValue(ctx, downloadBaseKey{}, newBase)
}

// Download downloads releases from go.dev/dl
func Download(ctx context.Context, all bool) (Releases, error) {
	return download(ctx, all)
}

func download(ctx context.Context, all bool) (Releases, error) {
	u, err := url.JoinPath(getBase(ctx), "/")
	if err != nil {
		return nil, err
	}
	u = u + "?mode=json"
	if all {
		u += "&include=all"
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200 (%d) for %s", res.StatusCode, u)
	}
	var rels Releases
	err = json.NewDecoder(res.Body).Decode(&rels)
	if err != nil {
		return nil, err
	}
	return rels, nil
}

// Download downloads a file as name.
// This fail if size or checksum are not match.
func (f File) Download(ctx context.Context, name string, force bool) error {
	if !force {
		ok, err := f.isDownloaded(name)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return f.download(ctx, name)
}

func (f File) download(ctx context.Context, name string) error {
	u, err := url.JoinPath(getBase(ctx), f.Filename)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code is not 200 (%d) for %s", res.StatusCode, u)
	}
	lf, err := os.Create(name)
	if err != nil {
		return err
	}
	h := sha256.New()
	r := io.TeeReader(res.Body, h)
	_, err = io.CopyN(lf, r, f.Size)
	lf.Close()
	if err != nil {
		os.Remove(name)
		return err
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if sum != f.ChecksumSHA256 {
		os.Remove(name)
		return fmt.Errorf("checksum mismatch expected=%s actual=%s for %s", f.ChecksumSHA256, sum, u)
	}
	return nil
}

func (f File) isDownloaded(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if fi.Size() != f.Size {
		return false, nil
	}
	r, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer r.Close()
	h := sha256.New()
	_, err = io.CopyN(h, r, f.Size)
	if err != nil {
		return false, err
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != f.ChecksumSHA256 {
		return false, nil
	}
	return true, nil
}
