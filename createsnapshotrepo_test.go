package artifactory_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ae6rt/artifactory"
)

func TestCreateSnapshot(t *testing.T) {
	doer := TestDoer{response: &http.Response{
		StatusCode: 200,
		Body:       nopCloser{bytes.NewBufferString("")},
	}}
	client := artifactory.NewClient(artifactory.Config{
		Username: "u",
		Password: "p",
		BaseURL:  "http://host:port",
		Doer:     &doer,
	})

	response, err := client.CreateSnapshotRepository("test-repo")
	if err != nil {
		t.Fatal(err)
	}
	if response != nil {
		t.Fatal("Unexpected response")
	}

	if doer.req.Method != "PUT" {
		t.Fatalf("Want PUT but got %s\n", doer.req.Method)
	}
	if doer.req.URL.Path != "/api/repositories/test-repo" {
		t.Fatalf("Want /api/repositories/test-repo but got %s\n", doer.req.URL)
	}
	if doer.req.Header.Get("Accept") != "*/*" {
		t.Fatalf("Want */* but got %s\n", doer.req.Header.Get("Accept"))
	}
	if doer.req.Header.Get("Content-type") != "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json" {
		t.Fatalf("want application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json but got %s\n", doer.req.Header.Get("Content-type"))
	}
	if doer.req.Header.Get("Authorization") == "" {
		t.Fatal("Want Basic Auth")
	}
}
