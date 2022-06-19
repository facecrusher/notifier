package processor

import (
	"testing"
	"time"

	"github.com/facecrusher/notifier/rest/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_MessageQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRestClient(ctrl)
	client.EXPECT().Post(gomock.Any(), gomock.Any()).Return(nil)

	testQueue := NewMessageQueue(nil, client)

	message := "this is a test"
	interval := time.Duration(time.Second)
	processed := make(chan string)
	testJob := NewNotificationJob(client, message, interval, &processed)

	testQueue.Start()

	// Send a job to the queue, wait for it to finish
	// and assert that processed message equals to sent message.
	// Also check that queue is empty (as message has been processed).
	done := make(chan bool)
	go func() {
		testQueue.internalQueue <- *testJob
		processedMsg := <-processed
		done <- true
		assert.Equal(t, testJob.Message.Message, processedMsg)
	}()
	<-done
	assert.True(t, len(testQueue.internalQueue) == 0)
}
