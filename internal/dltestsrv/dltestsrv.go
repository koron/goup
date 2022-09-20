package dltestsrv

import (
	_ "embed"
	"net/http"
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

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request) {
	// TODO:
}
