package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (c DefaultClient) CreateSnapshotRepository(repositoryID string) (*HTTPStatus, error) {
	repoConfig := LocalRepositoryConfiguration{
		Key:                     repositoryID,
		RClass:                  "local",
		Notes:                   "Created via automation with Artifactory Go client",
		PackageType:             "maven",
		RepoLayoutRef:           "maven-2-default",
		HandleSnapshots:         true,
		MaxUniqueSnapshots:      0,
		SnapshotVersionBehavior: "unique",
	}

	serial, err := json.Marshal(&repoConfig)
	if err != nil {
		return &HTTPStatus{}, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/repositories/%s", c.url, repositoryID), bytes.NewBuffer(serial))
	if err != nil {
		return &HTTPStatus{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")
	req.SetBasicAuth(c.user, c.password)

	response, err := c.client.Do(req)

	if err != nil {
		return &HTTPStatus{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return &HTTPStatus{StatusCode: response.StatusCode, Entity: data}, nil
	}

	return nil, nil
}

func NewClient(u, p, url string) Client {
	return DefaultClient{
		user:     u,
		password: p,
		url:      url,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
