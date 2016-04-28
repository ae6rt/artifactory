package artifactory

import "net/http"

type http500 struct {
	httpEntity []byte
}

type Client interface {
	CreateSnapshotRepository(string) (*HTTPStatus, error)
	GetVirtualRepositoryConfiguration(string) (VirtualRepositoryConfiguration, error)
	LocalRepositoryExists(string) (bool, error)
}

type HTTPStatus struct {
	StatusCode int
	Entity     []byte
}

type DefaultClient struct {
	user     string
	password string
	url      string
	client   *http.Client
}

type LocalRepositoryConfiguration struct {
	Key                     string `json:"key"`
	RClass                  string `json:"rclass"`
	Notes                   string `json:"notes"`
	PackageType             string `json:"packageType"`
	Description             string `json:"description"`
	RepoLayoutRef           string `json:"repoLayoutRef"`
	HandleSnapshots         bool   `json:"handleSnapshots"`
	MaxUniqueSnapshots      int    `json:"maxUniqueSnapshots"`
	SnapshotVersionBehavior string `json:"snapshotVersionBehavior"`
}

type VirtualRepositoryConfiguration struct {
	Key          string   `json:"key"`
	RClass       string   `json:"rclass"`
	Repositories []string `json:"repositories"`
	HTTPStatus   *HTTPStatus
}
