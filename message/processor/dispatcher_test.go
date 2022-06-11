package processor

import (
	"fmt"
	"notifier/message/domain"
	"notifier/rest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Dispatch(t *testing.T) {
	url := "http://www.test.com"
	headers := map[string]string{"Content-Type": "application/json"}

	// Dispatcher dependencies
	testClient := rest.NewNotifierRestClient(url, headers)
	testQueue := &MessageQueue{internalQueue: make(chan NotificationJob, 1)}
	testInterval := 1 * time.Second

	// New Dispatcher
	testDispatcher := NewDispatcher(*testQueue, *testClient, &testInterval)

	//Success case as messageQueue buffer is empty
	testMessage := domain.NewMessage("this is first test")
	err := testDispatcher.Dispatch(*testMessage)
	assert.Nil(t, err)

	//Fail case as messageQueue buffer is full
	testMessage = domain.NewMessage("this is second test")
	err = testDispatcher.Dispatch(*testMessage)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("processing queue full when retried. Message [%s] discarded", testMessage.Message), err.Error())
}

func Test_getInterval(t *testing.T) {
	assert.Equal(t, time.Duration(0), getInterval(nil))
	interval := 1 * time.Second
	assert.Equal(t, time.Second, getInterval(&interval))
}
