package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/facecrusher/notifier/message/domain"
	"github.com/facecrusher/notifier/rest/http_error"
	"github.com/facecrusher/notifier/rest/mock"
	"github.com/facecrusher/notifier/rest/response"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Post(t *testing.T) {
	ctrl := gomock.NewController(t)
	rb := mock.NewMockReqBuilder(ctrl)
	restClient := NewNotifierRestClient("http://www.test.com", nil)
	restClient.Request = rb
	defaultHeaders := make(http.Header)
	defaultHeaders.Add("Content-Type", "application/json")

	type testCase struct {
		name     string
		body     interface{}
		decode   interface{}
		response response.ReqResponse
		want     interface{}
	}
	testCases := []testCase{
		{
			name:   "response ok",
			body:   domain.Message{ID: "1234", Message: "this is a test message"},
			decode: make(map[string]string),
			response: response.ReqResponse{
				Response: &http.Response{
					StatusCode: http.StatusCreated,
					Header:     defaultHeaders,
				},
				ByteBody: []byte(`{}`),
			},
			want: nil,
		},
		{
			name:   "response with http error",
			body:   domain.Message{ID: "1234", Message: "this is a test message"},
			decode: make(map[string]string),
			response: response.ReqResponse{
				Response: &http.Response{
					StatusCode: http.StatusBadRequest,
					Request:    &http.Request{Header: defaultHeaders},
				},
				ByteBody: []byte(`{}`),
			},
			want: http_error.NewHTTPError("{}", http.StatusBadRequest, restClient.URL, defaultHeaders).Error(),
		},
		{
			name:   "response with unexpected error",
			body:   domain.Message{ID: "1234", Message: "this is a test message"},
			decode: make(map[string]string),
			response: response.ReqResponse{
				Err: errors.New("this is an error"),
			},
			want: errors.New("this is an error").Error(),
		},
		{
			name:   "response with decode error",
			body:   domain.Message{ID: "1234", Message: "this is a test message"},
			decode: make(map[string]string),
			response: response.ReqResponse{
				Response: &http.Response{
					StatusCode: http.StatusCreated,
					Header:     defaultHeaders,
				},
				ByteBody: []byte{},
			},
			want: "unexpected end of JSON input",
		},
	}
	for _, tc := range testCases {
		rb.EXPECT().DoPost(gomock.Any(), gomock.Any()).Return(&tc.response).AnyTimes()
		err := restClient.Post(tc.body, tc.decode)

		if err != nil {
			assert.NotNil(t, err)
			assert.Equal(t, tc.want, err.Error())
		} else {
			assert.Nil(t, err)
		}

	}
}
