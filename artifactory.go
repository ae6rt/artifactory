package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// NewClient returns a new Artifactory client with the given Config
func NewClient(config Config) Client {
	return DefaultClient{
		config: config,
	}
}

// CreateSnapshotRepository creates a snapshot repository with the given ID.  If the repository creation failed for reasons of transport failure,
// an error is returned.  If the repository creation failed for other business reasons, *HTTPStatus will have the details.
func (c DefaultClient) CreateSnapshotRepository(repositoryID string) (*HTTPStatus, error) {
	repoConfig := LocalRepositoryConfiguration{
		Key:                     repositoryID,
		RClass:                  "local",
		Notes:                   "Created via automation with https://github.com/ae6rt/artifactory Go client [" + time.Now().String() + "]",
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

	req, err := http.NewRequest("PUT", c.config.BaseURL+"/api/repositories/"+repositoryID, bytes.NewBuffer(serial))
	if err != nil {
		return &HTTPStatus{}, err
	}
	c.setAuthHeaders(req)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")

	response, err := c.config.Doer.Do(req)
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

// GetVirtualRepositoryConfiguration returns the configuration of the given virtual repository.
func (c DefaultClient) GetVirtualRepositoryConfiguration(repositoryID string) (VirtualRepositoryConfiguration, error) {

	req, err := http.NewRequest("GET", c.config.BaseURL+"/api/repositories/"+repositoryID, nil)
	if err != nil {
		return VirtualRepositoryConfiguration{}, err
	}
	c.setAuthHeaders(req)

	req.Header.Set("Accept", "application/vnd.org.jfrog.artifactory.repositories.VirtualRepositoryConfiguration+json")

	var data []byte
	var response *http.Response
	work := func() error {
		var err error
		response, err = c.config.Doer.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if data, err = ioutil.ReadAll(response.Body); err != nil {
			return err
		}
		return nil
	}
	err = retry(3, work)

	if err != nil {
		return VirtualRepositoryConfiguration{}, err
	}

	if response.StatusCode/100 == 5 {
		return VirtualRepositoryConfiguration{}, http500{data}
	}

	if response.StatusCode != 200 {
		return VirtualRepositoryConfiguration{HTTPStatus: &HTTPStatus{StatusCode: response.StatusCode, Entity: data}}, nil
	}

	var virtualRepository VirtualRepositoryConfiguration
	err = json.Unmarshal(data, &virtualRepository)
	return virtualRepository, err
}

// LocalRepositoryExists returns whether the given local repository exists.
func (c DefaultClient) LocalRepositoryExists(repositoryID string) (bool, error) {

	// https://www.jfrog.com/jira/browse/RTFACT-9998
	req, err := http.NewRequest("HEAD", c.config.BaseURL+"/api/repositories/"+repositoryID, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")
	c.setAuthHeaders(req)

	response, err := c.config.Doer.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	if response.StatusCode/100 == 5 {
		return false, http500{data}
	}

	return response.StatusCode == 200, nil
}

// RemoveRepository removes the given repository.  Check error for transport or marshaling errors.  Check HTTPStatus for other business errors.
func (c DefaultClient) RemoveRepository(repositoryID string) (*HTTPStatus, error) {
	req, err := http.NewRequest("DELETE", c.config.BaseURL+"/api/repositories/"+repositoryID, nil)
	if err != nil {
		return &HTTPStatus{}, err
	}
	c.setAuthHeaders(req)

	response, err := c.config.Doer.Do(req)
	if err != nil {
		return &HTTPStatus{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode/100 == 5 {
		return &HTTPStatus{StatusCode: response.StatusCode, Entity: data}, http500{data}
	}

	if response.StatusCode != 200 {
		return &HTTPStatus{response.StatusCode, data}, nil
	}

	return nil, nil
}

// RemoveItemFromRepository removes the given item from a repository.  Check error for transport or marshaling errors.
// Check HTTPStatus for other business errors.
func (c DefaultClient) RemoveItemFromRepository(repositoryID, item string) (*HTTPStatus, error) {
	if item == "" {
		panic("Refusing to remove an item of zero length.")
	}

	req, err := http.NewRequest("DELETE", c.config.BaseURL+"/api/repositories/"+repositoryID+"/"+item, nil)
	if err != nil {
		return &HTTPStatus{}, err
	}
	c.setAuthHeaders(req)

	response, err := c.config.Doer.Do(req)
	if err != nil {
		return &HTTPStatus{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode/100 == 5 {
		return &HTTPStatus{StatusCode: response.StatusCode, Entity: data}, http500{data}
	}

	if response.StatusCode != 200 {
		return &HTTPStatus{response.StatusCode, data}, nil
	}

	return nil, nil
}

// AddLocalRepositoryToGroup adds the given local repository to a virtual repository.  Check error for transport or marshaling errors.
// Check HTTPStatus for other business errors.
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

// RemoveLocalRepositoryFromGroup removes the given local repository to a virtual repository.  Check error for transport or marshaling errors.
// Check HTTPStatus for other business errors.
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

func remove(arr []string, removeIt string) []string {
	var t []string
	for _, v := range arr {
		if v == removeIt {
			continue
		}
		t = append(t, v)
	}
	return t
}

func (c DefaultClient) updateVirtualRepository(r VirtualRepositoryConfiguration) (*HTTPStatus, error) {
	serial, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.config.BaseURL+"/api/repositories/"+r.Key, bytes.NewBuffer(serial))
	if err != nil {
		return &HTTPStatus{}, err
	}
	c.setAuthHeaders(req)

	req.Header.Set("Accept", "*/*")
	// The Content-type prescribed by the API docs doesn't work:  https://www.jfrog.com/jira/browse/RTFACT-10035
	req.Header.Set("Content-type", "application/json")

	response, err := c.config.Doer.Do(req)
	if err != nil {
		return &HTTPStatus{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode/100 == 5 {
		return &HTTPStatus{StatusCode: response.StatusCode}, http500{data}
	}

	if response.StatusCode != 200 {
		return &HTTPStatus{response.StatusCode, data}, nil
	}
	return nil, nil
}

func (c DefaultClient) setAuthHeaders(req *http.Request) {
	if c.config.APIKey != "" {
		req.Header.Set("X-JFrog-Art-Api", c.config.APIKey)
	} else {
		req.SetBasicAuth(c.config.Username, c.config.Password)
	}
}

func retry(attempts int, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("retrying failed after %d attempts, last error: %s", attempts, err)
}
