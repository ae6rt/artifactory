package artifactory_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ae6rt/artifactory"
)

func TestRemoveSnapshot(t *testing.T) {
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

	response, err := client.RemoveRepository("test-repo")
	if err != nil {
		t.Fatal(err)
	}
	if response != nil {
		t.Fatal("Unexpected response")
	}

	if doer.req.Method != "DELETE" {
		t.Fatalf("wanted DELETE but found %s\n", doer.req.Method)
	}
	if doer.req.URL.Path != "/api/repositories/test-repo" {
		t.Fatalf("want /api/repositories/test-repo but found %s\n", doer.req.URL.Path)
	}
	if doer.req.Header.Get("Authorization") != "Basic dTpw" {
		t.Fatalf("Want  Basic dTpw but found %s\n", doer.req.Header.Get("Authorization"))
	}

}
