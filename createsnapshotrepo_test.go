package artifactory

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSnapshot(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("wanted PUT but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/api/repositories/test-repo" {
			t.Fatalf("want /api/repositories/test-repo but found %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("want application/json but got %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Content-type") != "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json" {
			t.Fatalf("want application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json but got %s\n", r.Header.Get("Content-type"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want  Basic dTpw but found %s\n", r.Header.Get("Authorization"))
		}
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	client := NewClient("u", "p", testServer.URL, &tls.Config{})
	_, err := client.CreateSnapshotRepository("test-repo")
	if err != nil {
		t.Fatal(err)
	}
}
