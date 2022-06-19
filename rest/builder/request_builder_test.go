package builder

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/facecrusher/notifier/rest/response"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func Test_DoPost(t *testing.T) {
	defaultHeaders := make(http.Header)
	defaultHeaders.Add("Content-Type", "application/json")
	defaultTimeout := time.Duration(time.Second)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK - Created"))
	}))

	rb := NewRequestBuilder(defaultHeaders, defaultTimeout)

	type testCase struct {
		name string
		url  string
		body interface{}
		want *response.ReqResponse
	}

	testCases := []testCase{
		{
			name: "DoPost ok request",
			url:  mockServer.URL,
			body: `{"message":"test message"}`,
			want: &response.ReqResponse{Response: &http.Response{StatusCode: http.StatusCreated}, ByteBody: []byte("OK - Created")},
		},
		{
			name: "DoPost error wrong URL format",
			url:  "test.com",
			body: "{}",
			want: &response.ReqResponse{Err: errors.New("Post \"test.com\": unsupported protocol scheme \"\"")},
		},
		{
			name: "DoPost error marshal request body",
			url:  "http://www.test.com",
			body: make(chan int),
			want: &response.ReqResponse{Err: errors.New("json: unsupported type: chan int")},
		},
	}

	for _, tc := range testCases {
		got := rb.DoPost(tc.url, tc.body)
		if tc.want.Err != nil {
			assert.Equal(t, tc.want.Err.Error(), got.Err.Error())
		} else {
			assert.NotNil(t, got)
			assert.Nil(t, got.Err)
			assert.Equal(t, tc.want.Response.StatusCode, got.StatusCode)
			assert.Equal(t, tc.want.ByteBody, got.ByteBody)
		}
	}
}
