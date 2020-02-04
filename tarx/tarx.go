package tarx

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DirInfo describes meta information of a dir.
type DirInfo struct {
	Name    string
	Mode    int64
	ModTime time.Time
}

// FileInfo describes meta information of a file.
type FileInfo struct {
	Name    string
	Size    int64
	Mode    int64
	ModTime time.Time
}

// Destination provides destination for extraction.
type Destination interface {
	// CreateDir creates a new directory in destination.
	CreateDir(info DirInfo) error

	// CreateFile creates a new file in destination.
	//
	// This can return io.WriteCloser as 1st return parameter, in that case
	// zipx close it automatically after have finished to use.
	CreateFile(info FileInfo) (io.Writer, error)
}

// ExtractFile extracts all files from a tar archive file "name".
func ExtractFile(ctx context.Context, name string, dst Destination) error {
	var uncompress func(io.Reader) (io.Reader, error)
	switch filepath.Ext(name) {
	case ".gz":
		uncompress = func(r io.Reader) (io.Reader, error) {
			return gzip.NewReader(r)
		}
	case ".bz2":
		uncompress = func(r io.Reader) (io.Reader, error) {
			return bzip2.NewReader(r), nil
		}
	case ".tar":
		uncompress = func(r io.Reader) (io.Reader, error) {
			return r, nil
		}
	default:
		return fmt.Errorf("unsupported archive: %s", name)
	}
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := uncompress(f)
	if err != nil {
		return err
	}
	err = Extract(ctx, r, dst)
	if err != nil {
		return err
	}
	return nil
}

// Extract extracts all files from `r` as a tar archive stream.
func Extract(ctx context.Context, r io.Reader, dst Destination) error {
	tr := tar.NewReader(r)
	for {
		err := ctx.Err()
		if err != nil {
			return err
		}
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch h.Typeflag {
		case tar.TypeDir:
			err := extractDir(dst, h)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := extractFile(dst, h, tr)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported tar's Typeflag: %c", h.Typeflag)
		}
	}
	return nil
}

func extractDir(dst Destination, h *tar.Header) error {
	err := dst.CreateDir(DirInfo{
		Name:    h.Name,
		Mode:    h.Mode,
		ModTime: h.ModTime,
	})
	if err != nil {
		return fmt.Errorf("failed to Destination.CreateDir name=%s: %w",
			h.Name, err)
	}
	return nil
}

func extractFile(dst Destination, h *tar.Header, tr *tar.Reader) error {
	w, err := dst.CreateFile(FileInfo{
		Name:    h.Name,
		Size:    h.Size,
		Mode:    h.Mode,
		ModTime: h.ModTime,
	})
	if err != nil {
		return fmt.Errorf("failed to Destination.CreateFile name=%s: %w", h.Name, err)
	}
	if wc, ok := w.(io.WriteCloser); ok {
		defer wc.Close()
	}
	_, err = io.CopyN(w, tr, h.Size)
	if err != nil {
		return err
	}
	return nil
}
