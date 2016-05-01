package artifactory_test

import (
	"io"
	"net/http"
)

type TestDoer struct {
	req      *http.Request
	response *http.Response
}

func (t *TestDoer) Do(req *http.Request) (*http.Response, error) {
	t.req = req
	return t.response, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
