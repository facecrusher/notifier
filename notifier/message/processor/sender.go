package processor

import (
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
)

type Sender interface {
	Send(message domain.Message) error
}

type MessageSender struct {
	RestClient rest.NotifierRestClient
}

func NewMessageSender(restClient rest.NotifierRestClient) *MessageSender {
	return &MessageSender{
		RestClient: restClient,
	}
}

func (ms *MessageSender) Send(message domain.Message) error {
	var decode map[string]string
	return ms.RestClient.Post(message, decode)
}
