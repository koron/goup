/*
Package dltestsrv provides test server for godlremote package.
*/
package dltestsrv

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

//go:embed active.json
var activeJSON []byte

//go:embed all.json
var allJSON []byte

//go:embed file.zip
var fileZip []byte

//go:embed file.tar.gz
var fileTarGz []byte

type Server struct {
	ReleaseActiveJSON []byte
	ReleaseAllJSON    []byte

	FileZip   []byte
	FileTarGz []byte
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

func (s *Server) fileZip() []byte {
	if s.FileZip != nil {
		return s.FileZip
	}
	return fileZip
}

func (s *Server) fileTarGz() []byte {
	if s.FileTarGz != nil {
		return s.FileTarGz
	}
	return fileTarGz
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
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", name))
		_, err = w.Write(s.fileZip())
	case "tar.gz":
		w.Header().Set("Content-Type", "application/x-gzip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar.gz\"", name))
		_, err = w.Write(s.fileTarGz())
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Printf("failed to write: %s", err)
	}
}
