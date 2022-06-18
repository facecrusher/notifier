package processor

import (
	"testing"
	"time"

	"github.com/facecrusher/notifier/rest/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRestClient(ctrl)
	client.EXPECT().Post(gomock.Any(), gomock.Any()).Return(nil)
	testMessage := "this is a test"
	interval := time.Duration(time.Second)

	testJob := NewNotificationJob(client, testMessage, interval)

	assert.Nil(t, testJob.Process())
}
