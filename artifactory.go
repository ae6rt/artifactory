// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Client API model based on:  https://github.com/google/go-github/blob/master/github/repos.go

package artifactory

import "fmt"

const (
	libraryVersion    = "0.1"
	userAgent         = "Xoom Artifactory Go SDK"
	defaultBaseURL    = "https://artifactory.example.com"
	defaultPathPrefix = "artifactory/"
)

type LocalRepositoryConfiguration struct {
	Key                     string `json:"key"`
	RClass                  string `json:"rclass"`
	Notes                   string `json:"notes"`
	PackageType             string `json:"packageType"`
	Description             string `json:"description"`
	RepoLayoutRef           string `json:"repoLayoutRef"`
	HandleSnapshots         bool   `json:"handleSnapshots"`
	HandleReleases          bool   `json:"handleReleases"`
	MaxUniqueSnapshots      int    `json:"maxUniqueSnapshots"`
	SnapshotVersionBehavior string `json:"snapshotVersionBehavior"`
}

type VirtualRepositoryConfiguration struct {
	Key           string   `json:"key"`
	RClass        string   `json:"rclass"`
	Repositories  []string `json:"repositories"`
	PackageType   string   `json:"packageType"`
	RepoLayoutRef string   `json:"repoLayoutRef"`
}

type RepositoryService struct {
	client     *Client
	PathPrefix string
}

// LocalConfiguration returns the configuration of the local repository named by repositoryKey.  The underlying HTTP Response is also
// returned, as well as any error.
func (service *RepositoryService) LocalConfiguration(repositoryKey string) (*LocalRepositoryConfiguration, *Response, error) {
	u := fmt.Sprintf("%sapi/repositories/%s", service.PathPrefix, repositoryKey)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	p := new(LocalRepositoryConfiguration)
	resp, err := service.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, err
}

// VirtualConfiguration returns the configuration for the virtual repository named by repositoryKey.  The underlying HTTP Response is also
// returned, as well as any error.
func (service *RepositoryService) VirtualConfiguration(repositoryKey string) (*VirtualRepositoryConfiguration, *Response, error) {
	u := fmt.Sprintf("%sapi/repositories/%s", service.PathPrefix, repositoryKey)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	p := new(VirtualRepositoryConfiguration)
	resp, err := service.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, err
}

// Remove removes the repository named by repositoryKey, and returns the underlying HTTP response and any error.
func (service *RepositoryService) Remove(repositoryKey string) (*Response, error) {
	u := fmt.Sprintf("%sapi/repositories/%s", service.PathPrefix, repositoryKey)
	req, err := service.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// AddToGroup adds the local repository named by localRepositoryID to the virtual repository named by virtualRepositoryID.  The underlying HTTP response is also returned,
// as well as any error.
func (service *RepositoryService) AddToGroup(virtualRepositoryID, localRepositoryID string) (*Response, error) {
	virtual, response, err := service.VirtualConfiguration(virtualRepositoryID)
	if err != nil {
		return response, err
	}

	if (response.StatusCode / 100) != 2 {
		return response, nil
	}

	if contains(virtual.Repositories, localRepositoryID) {
		return nil, nil
	}

	virtual.Repositories = append(virtual.Repositories, localRepositoryID)

	return service.updateVirtualRepository(virtual)
}

func (service *RepositoryService) updateVirtualRepository(virtualRepository *VirtualRepositoryConfiguration) (*Response, error) {
	panic("NYI")
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
