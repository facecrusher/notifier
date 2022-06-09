package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RequestBuilder struct {
	Headers http.Header
	Timeout time.Duration
	Client  *http.Client
}

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

func (rb *RequestBuilder) getClient() http.Client {
	client := &http.Client{}
	client.Timeout = rb.Timeout
	return *client
}

func parseURL(reqlURL string) (string, error) {
	rURL, err := url.Parse(reqlURL)
	if err != nil {
		return reqlURL, err
	}
	return rURL.String(), nil
}

func marshalReqBody(body interface{}) (b []byte, err error) {
	if body != nil {
		b, err = json.Marshal(body)
	}
	return
}
