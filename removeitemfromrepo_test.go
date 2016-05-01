package artifactory_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ae6rt/artifactory"
)

func TestRemoveItemFromRepo(t *testing.T) {
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

	response, err := client.RemoveItemFromRepository("test-repo", "item")
	if err != nil {
		t.Fatal(err)
	}
	if response != nil {
		t.Fatal("Unexpected response")
	}
	if err != nil {
		t.Fatal(err)
	}

	if doer.req.Method != "DELETE" {
		t.Fatalf("Want DELETE but got %s\n", doer.req.Method)
	}
	if doer.req.URL.Path != "/api/repositories/test-repo/item" {
		t.Fatalf("Want /api/repositories/test-repo/item but got %s\n", doer.req.URL)
	}
	if doer.req.Header.Get("Authorization") == "" {
		t.Fatal("Want Basic Auth")
	}
}
