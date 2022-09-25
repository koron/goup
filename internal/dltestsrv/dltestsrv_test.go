package dltestsrv_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/koron/goup/internal/dltestsrv"
)

func TestActiveReleases(t *testing.T) {
	srv := httptest.NewServer(&dltestsrv.Server{})
	defer srv.Close()
	r, err := srv.Client().Get(srv.URL + "/?mode=json")
	if err != nil {
		t.Fatal(err)
	}
	var rels []map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&rels)
	r.Body.Close()
	if err != nil {
		t.Fatalf("invalid JSON in a response: %s", err)
	}

	if len(rels) != 2 {
		t.Fatalf("want 2 releases but got %d", len(rels))
	}
	r0 := rels[0]
	if r0["version"] != "go1.19.1" {
		t.Errorf("[0].version want \"go1.19.1\" but got %q", r0["version"])
	}
	if r0["stable"] != true {
		t.Errorf("[0].stable want true but got %t", r0["stable"])
	}
	r1 := rels[1]
	if r1["version"] != "go1.18.6" {
		t.Errorf("[1].version want \"go1.18.6\" but got %q", r1["version"])
	}
	if r1["stable"] != true {
		t.Errorf("[1].stable want true but got %t", r1["stable"])
	}
}

func TestAllReleases(t *testing.T) {
	srv := httptest.NewServer(&dltestsrv.Server{})
	defer srv.Close()
	r, err := srv.Client().Get(srv.URL + "/?mode=json&include=all")
	if err != nil {
		t.Fatal(err)
	}
	var rels []map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&rels)
	r.Body.Close()
	if err != nil {
		t.Fatalf("invalid JSON in a response: %s", err)
	}
	if len(rels) != 232 {
		t.Fatalf("want 232 releases but got %d", len(rels))
	}
}

func TestFileZip(t *testing.T) {
	srv := httptest.NewServer(&dltestsrv.Server{})
	defer srv.Close()
	r, err := srv.Client().Get(srv.URL + "/go1.19.1.windows-amd64.zip")
	if err != nil {
		t.Fatal(err)
	}
	r.Body.Close()
	// TODO: test r.Body as zip
}

func TestFileTarGz(t *testing.T) {
	srv := httptest.NewServer(&dltestsrv.Server{})
	defer srv.Close()
	r, err := srv.Client().Get(srv.URL + "/go1.19.1.linux-amd64.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	r.Body.Close()
	// TODO: test r.Body as tar.gz
}
