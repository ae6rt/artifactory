package artifactory

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoveItemFromRepo(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("wanted DELETE but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/api/repositories/test-repo/item" {
			t.Fatalf("want /api/repositories/test-repo/item but found %s\n", url.Path)
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want  Basic dTpw but found %s\n", r.Header.Get("Authorization"))
		}
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	client := NewBasicAuthClient("u", "p", testServer.URL, &tls.Config{})
	_, err := client.RemoveItemFromRepository("test-repo", "item")
	if err != nil {
		t.Fatal(err)
	}
}
