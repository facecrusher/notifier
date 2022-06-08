package processor

import (
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SendMessage(t *testing.T) {
	mockMessage := domain.NewMessage("This is a test")
	url := "https://eouss1txxyn5t7x.m.pipedream.net"
	headers := make(map[string]string)
	client := rest.NewNotifierRestClient(url, headers)
	readyPool := make(chan chan domain.Message)
	done := sync.WaitGroup{}
	messageSender := NewMessageSender(*client, readyPool, done)

	err := messageSender.Send(*mockMessage)
	assert.Nil(t, err)
}
