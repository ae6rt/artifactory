package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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

// GetVirtualRepositoryConfiguration retrieves virtual repository configuration.  Whether an error is returned or not
// is driven by whether a retry framework shuuld retry such a call.
func (c DefaultClient) GetVirtualRepositoryConfiguration(repositoryID string) (VirtualRepositoryConfiguration, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/repositories/%s", c.url, repositoryID), nil)
	if err != nil {
		return VirtualRepositoryConfiguration{}, err
	}

	req.Header.Set("Accept", "application/vnd.org.jfrog.artifactory.repositories.VirtualRepositoryConfiguration+json")
	req.SetBasicAuth(c.user, c.password)

	response, err := c.client.Do(req)
	if err != nil {
		return VirtualRepositoryConfiguration{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return VirtualRepositoryConfiguration{}, err
	}

	if response.StatusCode/100 == 5 {
		return VirtualRepositoryConfiguration{}, http500{data}
	}

	var virtualRepository VirtualRepositoryConfiguration
	if err := json.Unmarshal(data, &virtualRepository); err != nil {
		return VirtualRepositoryConfiguration{}, err
	}

	if response.StatusCode != 200 {
		return VirtualRepositoryConfiguration{HTTPStatus: &HTTPStatus{StatusCode: response.StatusCode, Entity: data}}, nil
	}

	return virtualRepository, nil
}

func (c DefaultClient) LocalRepositoryExists(repositoryID string) (bool, error) {

	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/api/repositories/%s", c.url, repositoryID), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")
	req.SetBasicAuth(c.user, c.password)

	response, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode/100 == 5 {
		return false, http500{}
	}

	if response.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}

func (h http500) Error() string {
	return string(h.httpEntity)
}
