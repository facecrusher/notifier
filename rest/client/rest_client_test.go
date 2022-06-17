package client

import (
	"errors"
	"testing"

	"github.com/facecrusher/notifier/message/domain"
	"github.com/facecrusher/notifier/rest/mock"
	"github.com/facecrusher/notifier/rest/response"
	"github.com/golang/mock/gomock"
)

func Test_Post(t *testing.T) {
	ctrl := gomock.NewController(t)
	rb := mock.NewMockReqBuilder(ctrl)

	type testCase struct {
		name   string
		body   interface{}
		decode interface{}
		want   response.ReqResponse
	}
	testCases := []testCase{
		{
			name:   "response with error",
			body:   domain.Message{ID: "1234", Message: "this is a test message"},
			decode: make(map[string]string),
			want:   response.ReqResponse{Err: errors.New("this is an error")},
		},
	}
	for _, tc := range testCases {
		rb.EXPECT().DoPost(gomock.Any(), gomock.Any()).Return(tc.want)

	}
}
