// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Client API model based on:  https://github.com/google/go-github/blob/master/github/repos.go

package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
   An ErrorResponse reports one or more errors caused by an API request.

   GitHub API docs: http://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
	// Block is only populated on certain types of errors such as code 451.
	// See https://developer.github.com/changes/2016-03-17-the-451-status-code-is-now-supported/
	// for more information.
	Block *struct {
		Reason string `json:"reason,omitempty"`
		//CreatedAt *Timestamp `json:"created_at,omitempty"`
	} `json:"block,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Errors)
}

// An Error reports more details on an individual error in an ErrorResponse.
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
}

type Response struct {
	*http.Response

	//NextPage  int
	//PrevPage  int
	////FirstPage int
	//LastPage  int

	//Rate
}

type Client struct {
	client            *http.Client
	BaseURL           *url.URL
	RepositoryService *RepositoryService
	UserAgent         string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	b, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: b, UserAgent: userAgent}
	c.RepositoryService = &RepositoryService{client: c}
	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("@@@ %s base url\n", c.BaseURL.String())

	u := c.BaseURL.ResolveReference(rel)
	fmt.Printf("@@@ %s resolved url\n", u)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.  If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		//	io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	//response.populatePageValues()
	return response
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	switch {
	//case r.StatusCode == http.StatusUnauthorized && strings.HasPrefix(r.Header.Get(headerOTP), "required"):
	/*
		        return (*TwoFactorAuthError)(errorResponse)
				            case r.StatusCode == http.StatusForbidden && r.Header.Get(headerRateRemaining) == "0" && strings.HasPrefix(errorResponse.Message, "API rate limit exceeded for "):
							                return &RateLimitError{
											                    Rate:     parseRate(r),
																                    Response: errorResponse.Response,
																					                    Message:  errorResponse.Message,
																										                }
	*/
	default:
		return errorResponse
	}
}
