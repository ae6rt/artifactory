package artifactory

import (
	"net/http"
	"testing"
)

func TestSetAPIKeyAuthHeaders(t *testing.T) {
	client := DefaultClient{Config{APIKey: "abc"}}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	client.setAuthHeaders(req)
	if req.Header.Get("X-JFrog-Art-Api") != "abc" {
		t.Fatalf("Want abc but got %s\n", req.Header.Get("J-Frog-Art-Api"))
	}
}

func TestSetBasicAuthHeaders(t *testing.T) {
	client := DefaultClient{Config{Username: "bob", Password: "sekrit"}}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	client.setAuthHeaders(req)
	if req.Header.Get("Authorization") == "" {
		t.Fatal("Expecting Basic Auth header but found none")
	}
}
