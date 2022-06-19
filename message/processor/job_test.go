package processor

import (
	"errors"
	"testing"
	"time"

	"github.com/facecrusher/notifier/rest/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRestClient(ctrl)
	testMessage := "this is a test"
	interval := time.Duration(time.Second)

	testJob := NewNotificationJob(client, testMessage, interval, nil)

	client.EXPECT().Post(gomock.Any(), gomock.Any()).Return(nil)
	assert.Nil(t, testJob.Process())

	client.EXPECT().Post(gomock.Any(), gomock.Any()).Return(errors.New("process error"))
	err := testJob.Process()
	assert.Error(t, err)
	assert.Equal(t, "process error", err.Error())

}
