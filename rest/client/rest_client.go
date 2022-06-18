package client

import (
	"net/http"
	"time"

	"github.com/facecrusher/notifier/rest/builder"
	"github.com/facecrusher/notifier/rest/http_error"
)

const (
	DEFAULT_TIMEOUT = 1000 * time.Millisecond
)

//go:generate mockgen -source=./rest_client.go -destination=./mock/rest_client.go -package=mock
type RestClient interface {
	Post(body interface{}, decode interface{}) error
}

type NotifierRestClient struct {
	Request builder.ReqBuilder
	URL     string
}

func NewNotifierRestClient(url string, headers *map[string]string) *NotifierRestClient {
	reqHeaders := setHeaders(headers)
	return &NotifierRestClient{
		URL:     url,
		Request: builder.NewRequestBuilder(reqHeaders, DEFAULT_TIMEOUT),
	}
}

func (nrc *NotifierRestClient) Post(body interface{}, decode interface{}) error {
	response := nrc.Request.DoPost(nrc.URL, body)
	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode != http.StatusCreated {
		cause := string(response.Bytes())
		return http_error.NewHTTPError(cause, response.StatusCode, nrc.URL, response.Request.Header)
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
