package builder

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/facecrusher/notifier/rest/response"
)

//go:generate mockgen -source=./request_builder.go -destination=./mock/request_builder.go -package=mock
type ReqBuilder interface {
	DoPost(url string, body interface{}) (responseObj *response.ReqResponse)
}

type RequestBuilder struct {
	Headers http.Header
	Timeout time.Duration
	Client  *http.Client
}

func NewRequestBuilder(headers http.Header, timeout time.Duration) ReqBuilder {
	return &RequestBuilder{
		Headers: headers,
		Timeout: timeout,
	}
}

// DoPost builds and executes a post request to the given url with the corresponding body as payload
func (rb *RequestBuilder) DoPost(url string, body interface{}) (responseObj *response.ReqResponse) {
	responseObj = new(response.ReqResponse)
	client := rb.getClient()

	// Parse URL
	resourceURL, err := parseURL(url)
	if err != nil {
		responseObj.Err = err
		return
	}

	// Marshal request to JSON
	reqBody, err := marshalReqBody(body)
	if err != nil {
		responseObj.Err = err
		return
	}

	// Create request object
	req, err := http.NewRequest(http.MethodPost, resourceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		responseObj.Err = err
		return
	}
	// Set headers if any defined
	req.Header = rb.Headers

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		responseObj.Err = err
		return
	}
	// Read response
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responseObj.Err = err
		return
	}
	responseObj.Response = resp
	responseObj.ByteBody = respBody
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
