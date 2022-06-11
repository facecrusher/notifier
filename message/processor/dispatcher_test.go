package processor

import (
	"notifier/message/domain"
	"notifier/rest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Dispatch(t *testing.T) {
	url := "http://www.test.com"
	headers := map[string]string{"Content-Type": "application/json"}
	client := rest.NewNotifierRestClient(url, headers)
	testDispatcher := &MessageDispatcher{
		MessageQueue: make(chan NotificationJob),
		Client:       *client,
		Interval:     1 * time.Second,
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
