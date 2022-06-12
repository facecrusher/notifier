package rest

import (
	"net/http"
	"time"
)

const (
	DEFAULT_TIMEOUT = 1000 * time.Millisecond
)

type RestClient interface {
	Post(body interface{}, decode interface{}) error
}

type NotifierRestClient struct {
	Request *RequestBuilder
	URL     string
}

func NewNotifierRestClient(url string, headers *map[string]string) *NotifierRestClient {
	reqHeaders := setHeaders(headers)
	return &NotifierRestClient{
		URL: url,
		Request: &RequestBuilder{
			Headers: reqHeaders,
			Timeout: DEFAULT_TIMEOUT,
		},
	}
}

func (nrc *NotifierRestClient) Post(body interface{}, decode interface{}) error {
	response := nrc.Request.DoPost(nrc.URL, body)
	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode != http.StatusOK {
		cause := string(response.Bytes())
		return NewHTTPError(cause, response.StatusCode, nrc.URL, response.Request.Header)
	}

	err := response.Decode(&decode)
	if err != nil {
		return err
	}
	return nil
}

func setHeaders(headers *map[string]string) http.Header {
	reqHeaders := make(http.Header)
	reqHeaders.Add("Content-Type", "application/json")
	if headers != nil {
		for k, v := range *headers {
			reqHeaders.Add(k, v)
		}
	}
	return reqHeaders
}
