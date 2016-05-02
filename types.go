package artifactory

import (
	"log"
	"net/http"
)

type http500 struct {
	httpEntity []byte
}

type Config struct {
	Username string
	Password string
	APIKey   string
	BaseURL  string
	Doer     Doer
	Log      *log.Logger
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type DefaultClient struct {
	config Config
}

type Client interface {
	CreateSnapshotRepository(string) (*HTTPStatus, error)
	RemoveRepository(string) (*HTTPStatus, error)
	LocalRepositoryExists(string) (bool, error)
	GetVirtualRepositoryConfiguration(string) (VirtualRepositoryConfiguration, error)
	AddLocalRepositoryToGroup(string, string) (*HTTPStatus, error)
	RemoveLocalRepositoryFromGroup(string, string) (*HTTPStatus, error)
	RemoveItemFromRepository(string, string) (*HTTPStatus, error)
}

type HTTPStatus struct {
	StatusCode int
	Entity     []byte
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
