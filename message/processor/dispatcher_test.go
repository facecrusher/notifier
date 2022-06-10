package processor

import (
	"notifier/message/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dispatch(t *testing.T) {
	testDispatcher := &MessageDispatcher{
		MessageQueue: make(chan NotificationJob),
	}

	testMessage := domain.NewMessage("this is a test")
	testDispatcher.Dispatch(*testMessage)

}

func Test_IsQueueClosed(t *testing.T) {
	testChan := make(chan NotificationJob)
	assert.False(t, isQueueClosed(testChan))

	close(testChan)
	assert.True(t, isQueueClosed(testChan))
}
