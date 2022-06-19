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
	defer testQueue.Stop()

	done := make(chan bool)

	// Test Process Message
	go func() {
		testQueue.internalQueue <- *testJob //Send test job to the queue
		processedMsg := <-processed         //Wait for it to be processed
		done <- true
		assert.Equal(t, testJob.Message.Message, processedMsg) //Assert processed msg is same as test job message
	}()
	<-done
	assert.True(t, len(testQueue.internalQueue) == 0) // Assert queue is empty after message is processed
}
