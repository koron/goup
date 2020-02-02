package godlremote

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const dlURL = "https://golang.org/dl/"

// Download downloads releases from golang.org/dl
func Download(ctx context.Context, all bool) (Releases, error) {
	return download(ctx, dlURL, all)
}

func download(ctx context.Context, base string, all bool) (Releases, error) {
	u := dlURL + "?mode=json"
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
func (f File) Download(ctx context.Context, name string) error {
	u := dlURL + f.Filename
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
