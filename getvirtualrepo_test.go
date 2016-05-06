package artifactory_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ae6rt/artifactory"
)

func TestGetVirtualRepositoryConfiguration(t *testing.T) {
	json := `
{
  "key" : "test-repo",
  "packageType" : "maven",
  "description" : "automation",
  "notes" : "",
  "includesPattern" : "**/*",
  "excludesPattern" : "",
  "repoLayoutRef" : "maven-2-default",
  "rclass" : "virtual"
}
`
	doer := TestDoer{response: &http.Response{
		StatusCode: 200,
		Body:       nopCloser{bytes.NewBufferString(json)},
	}}
	client := artifactory.NewClient(artifactory.Config{
		Username: "u",
		Password: "p",
		BaseURL:  "http://host:port",
		Doer:     &doer,
	})

	response, err := client.GetVirtualRepositoryConfiguration("test-repo")
	if err != nil {
		t.Fatal(err)
	}
	if response.HTTPStatus != nil {
		t.Fatal("Unexpected response")
	}
	if response.Key != "test-repo" {
		t.Fatalf("want test-repo but got %s\n", response.Key)
	}
	if response.RClass != "virtual" {
		t.Fatalf("want virtual but got %s\n", response.RClass)
	}

	if doer.req.Method != "GET" {
		t.Fatalf("wanted GET but found %s\n", doer.req.Method)
	}
	if doer.req.URL.Path != "/api/repositories/test-repo" {
		t.Fatalf("want /api/repositories/test-repo but found %s\n", doer.req.URL.Path)
	}
	if doer.req.Header.Get("Authorization") != "Basic dTpw" {
		t.Fatalf("Want  Basic dTpw but found %s\n", doer.req.Header.Get("Authorization"))
	}

}
