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
