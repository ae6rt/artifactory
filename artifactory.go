package artifactory

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func NewApiKeyClient(apiKey, url string, tlsConfig *tls.Config) Client {
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return DefaultClient{
		apiKey: apiKey,
		url:    url,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

func NewBasicAuthClient(username, password, url string, tlsConfig *tls.Config) Client {
	return DefaultClient{
		user:     username,
		password: password,
		url:      url,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
	}
}

func (c DefaultClient) CreateSnapshotRepository(repositoryID string) (*HTTPStatus, error) {
	repoConfig := LocalRepositoryConfiguration{
		Key:                     repositoryID,
		RClass:                  "local",
		Notes:                   "Created via automation with Artifactory Go client [" + time.Now().String() + "]",
		PackageType:             "maven",
		RepoLayoutRef:           "maven-2-default",
		HandleSnapshots:         true,
		HandleReleases:          false,
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

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")
	if c.apiKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.apiKey)
	} else {
		req.SetBasicAuth(c.user, c.password)
	}

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
	if c.apiKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.apiKey)
	} else {
		req.SetBasicAuth(c.user, c.password)
	}

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
	if c.apiKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.apiKey)
	} else {
		req.SetBasicAuth(c.user, c.password)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode/100 == 5 {
		return false, http500{}
	}

	return response.StatusCode == 200, nil
}

func (c DefaultClient) RemoveSnapshotRepository(repositoryID string) (*HTTPStatus, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/repositories/%s", c.url, repositoryID), nil)
	if err != nil {
		return &HTTPStatus{}, err
	}

	if c.apiKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.apiKey)
	} else {
		req.SetBasicAuth(c.user, c.password)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return &HTTPStatus{}, err
	}
	defer response.Body.Close()

	if response.StatusCode/100 == 5 {
		return &HTTPStatus{}, http500{}
	}

	return nil, nil
}

func (c DefaultClient) AddLocalRepositoryToGroup(virtualRepositoryID, localRepositoryID string) (*HTTPStatus, error) {
	r, err := c.GetVirtualRepositoryConfiguration(virtualRepositoryID)
	if err != nil {
		return nil, err
	}
	if r.HTTPStatus != nil {
		return r.HTTPStatus, nil
	}

	if contains(r.Repositories, localRepositoryID) {
		return nil, nil
	}

	r.Repositories = append(r.Repositories, localRepositoryID)

	return c.updateVirtualRepository(r)
}

func (c DefaultClient) RemoveLocalRepositoryFromGroup(virtualRepositoryID, localRepositoryID string) (*HTTPStatus, error) {
	r, err := c.GetVirtualRepositoryConfiguration(virtualRepositoryID)
	if err != nil {
		return nil, err
	}
	if r.HTTPStatus != nil {
		return r.HTTPStatus, nil
	}

	if !contains(r.Repositories, localRepositoryID) {
		return nil, nil
	}

	r.Repositories = remove(r.Repositories, localRepositoryID)

	return c.updateVirtualRepository(r)
}

func (h http500) Error() string {
	return string(h.httpEntity)
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func remove(arr []string, value string) []string {
	var t []string
	for _, v := range arr {
		if v != value {
			t = append(t, value)
		}
	}
	return t
}

func (c DefaultClient) updateVirtualRepository(r VirtualRepositoryConfiguration) (*HTTPStatus, error) {
	serial, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/repositories/%s", c.url, r.Key), bytes.NewBuffer(serial))
	if err != nil {
		return &HTTPStatus{}, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/vnd.org.jfrog.artifactory.repositories.VirtualRepositoryConfiguration+json")
	if c.apiKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.apiKey)
	} else {
		req.SetBasicAuth(c.user, c.password)
	}

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
		return &HTTPStatus{response.StatusCode, data}, nil
	}
	return nil, nil
}
