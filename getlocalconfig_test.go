package artifactory

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestGetLocalConfiguration(t *testing.T) {
	user := os.Getenv("TEST_USER")
	password := os.Getenv("TEST_PASSWORD")
	testURL := os.Getenv("TEST_URL")
	bt := &BasicAuthTransport{Username: user, Password: password}
	c := bt.Client()

	client := NewClient(c)
	client.BaseURL, _ = url.Parse(testURL)

	config, response, err := client.RepositoryService.LocalConfiguration("inftools-local")
	fmt.Printf("%+v, %+v, %+v\n", config, response, err)
}
