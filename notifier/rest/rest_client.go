package rest

import (
	"net/http"
)

type RestClient interface {
	Post(url string, body interface{}, decode interface{}, header http.Header) error
}

type NotifierRestClient struct {
	Request *RequestBuilder
	URL     string
}

func NewNotifierRestClient(url string, headers map[string]string) *NotifierRestClient {
	defaultHeaders := make(http.Header)
	for k, v := range headers {
		defaultHeaders.Add(k, v)
	}
	return &NotifierRestClient{
		URL: url,
		Request: &RequestBuilder{
			Headers: defaultHeaders,
			Timeout: 10,
		},
	}
}

func (nrc *NotifierRestClient) Post(body interface{},
	decode interface{}) error {
	response := nrc.Request.DoPost(nrc.URL, body)
	if response.Err != nil {
		return nil
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
