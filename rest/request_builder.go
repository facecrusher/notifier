package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ReqBuilder interface {
	DoPost(url string, body interface{}) (response *Response)
}

type RequestBuilder struct {
	Headers http.Header
	Timeout time.Duration
	Client  *http.Client
}

// DoPost builds and executes a post request to the given url with the corresponding body as payload
func (rb *RequestBuilder) DoPost(url string, body interface{}) (response *Response) {
	response = new(Response)
	client := rb.getClient()
	// Parse URL
	resourceURL, err := parseURL(url)
	if err != nil {
		response.Err = err
		return
	}

	// Marshal request to JSON
	reqBody, err := marshalReqBody(body)
	if err != nil {
		response.Err = err
		return
	}

	// Execute POST request
	resp, err := client.Post(resourceURL, "appplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		response.Err = err
		return
	}

	// Read response
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Err = err
		return
	}

	response.Response = resp
	response.byteBody = respBody
	return
}

// getClient returns a configurable http client to perform requests
func (rb *RequestBuilder) getClient() http.Client {
	client := &http.Client{}
	client.Timeout = rb.Timeout
	return *client
}

// parseURL parses reqURL into a URL structure
func parseURL(reqlURL string) (string, error) {
	rURL, err := url.Parse(reqlURL)
	if err != nil {
		return reqlURL, err
	}
	return rURL.String(), nil
}

// marshalReqBody converts the body interface into a byte array for processing
func marshalReqBody(body interface{}) (b []byte, err error) {
	if body != nil {
		b, err = json.Marshal(body)
	}
	return
}
