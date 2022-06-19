package builder

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/facecrusher/notifier/rest/response"
	"github.com/stretchr/testify/assert"
)

func Test_DoPost(t *testing.T) {
	defaultHeaders := make(http.Header)
	defaultHeaders.Add("Content-Type", "application/json")
	defaultTimeout := time.Duration(time.Second)
	rb := NewRequestBuilder(defaultHeaders, defaultTimeout)

	type testCase struct {
		name string
		url  string
		body interface{}
		want *response.ReqResponse
	}

	testCases := []testCase{
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
		if got.Err != nil {
			assert.Equal(t, tc.want.Err.Error(), got.Err.Error())
		}
	}
}
