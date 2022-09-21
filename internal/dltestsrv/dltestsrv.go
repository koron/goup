package dltestsrv

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	_ "embed"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

//go:embed active.json
var activeJSON []byte

//go:embed all.json
var allJSON []byte

type Server struct {
	ReleaseActiveJSON []byte
	ReleaseAllJSON    []byte
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.URL.Path == "/" {
		s.serveReleases(w, r)
		return
	}
	s.serveFile(w, r)
}

func (s *Server) serveReleases(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("mode") != "json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if q.Get("include") == "all" {
		w.Write(s.releaseAll())
		return
	}
	w.Write(s.releaseActive())
}

func (s *Server) releaseActive() []byte {
	if s.ReleaseActiveJSON != nil {
		return s.ReleaseActiveJSON
	}
	return activeJSON
}

func (s *Server) releaseAll() []byte {
	if s.ReleaseAllJSON != nil {
		return s.ReleaseAllJSON
	}
	return allJSON
}

var rxGoFile = regexp.MustCompile(`\b(go\d+(?:\.\d+)*(?:(?:rc|beta|alpha)\d+)?\.(?:\D[^-]*)-(?:[^.]+))\.(tar\.gz|zip)$`)

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request) {
	m := rxGoFile.FindStringSubmatch(r.URL.Path)
	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name, ext := m[1], m[2]
	var err error
	switch ext {
	case "zip":
		w.Header().Set("Content-Type", "application/zip")
		err = s.writeFileZip(w, name)
	case "tar.gz":
		w.Header().Set("Content-Type", "application/x-gzip")
		err = s.writeFileTarGz(w, name)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Printf("failed to write: %w", err)
	}
}

func (s *Server) writeFileZip(w io.Writer, name string) error {
	zw := zip.NewWriter(w)
	defer zw.Close()
	tw, err := zw.CreateHeader(&zip.FileHeader{
		Name:     "go/README.txt",
		Method:   zip.Deflate,
		Modified: time.Now(),
	})
	if err != nil {
		return err
	}
	_, err = io.WriteString(tw, name+"\n")
	return err
}

func (s *Server) writeFileTarGz(w io.Writer, name string) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	readmeBody := name + "\n"
	err := tw.WriteHeader(&tar.Header{
		Name:    "go/README.txt",
		Mode:    0644,
		Size:    int64(len(readmeBody)),
		Uname:   "root",
		Gname:   "root",
		ModTime: time.Now(),
	})
	if err != nil {
		return err
	}
	io.WriteString(tw, readmeBody)
	tw.Flush()
	return nil
}
