package artifactory

import "net/http"

type http500 struct {
	httpEntity []byte
}

type Client interface {
	CreateSnapshotRepository(string) (*HTTPStatus, error)
	LocalRepositoryExists(string) (bool, error)
	GetVirtualRepositoryConfiguration(string) (VirtualRepositoryConfiguration, error)
	AddLocalRepositoryToGroup(string, string) (*HTTPStatus, error)
	RemoveLocalRepositoryFromGroup(string, string) (*HTTPStatus, error)
}

type HTTPStatus struct {
	StatusCode int
	Entity     []byte
}

type DefaultClient struct {
	user     string
	password string
	apiKey   string
	url      string
	client   *http.Client
}

type LocalRepositoryConfiguration struct {
	Key                     string      `json:"key"`
	RClass                  string      `json:"rclass"`
	Notes                   string      `json:"notes"`
	PackageType             string      `json:"packageType"`
	Description             string      `json:"description"`
	RepoLayoutRef           string      `json:"repoLayoutRef"`
	HandleSnapshots         bool        `json:"handleSnapshots"`
	HandleReleases          bool        `json:"handleReleases"`
	MaxUniqueSnapshots      int         `json:"maxUniqueSnapshots"`
	SnapshotVersionBehavior string      `json:"snapshotVersionBehavior"`
	HTTPStatus              *HTTPStatus `json:"-"`
}

type VirtualRepositoryConfiguration struct {
	Key           string      `json:"key"`
	RClass        string      `json:"rclass"`
	Repositories  []string    `json:"repositories"`
	PackageType   string      `json:"packageType"`
	RepoLayoutRef string      `json:"repoLayoutRef"`
	HTTPStatus    *HTTPStatus `json:"-"`
}

type BooleanResponse struct {
	Result     bool
	HTTPStatus *HTTPStatus
}
