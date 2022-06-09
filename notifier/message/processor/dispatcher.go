package processor

import (
	"errors"
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
)

type Dispatcher interface {
	Dispatch(message domain.Message) error
}

type MessageDispatcher struct {
	Client       rest.NotifierRestClient
	MessageQueue chan NotificationJob
	IsFinished   bool
}

func NewDispatcher(mq MessageQueue, c rest.NotifierRestClient) *MessageDispatcher {
	return &MessageDispatcher{
		Client:       c,
		MessageQueue: mq.GetMessageQueue(),
		IsFinished:   false,
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	if md.IsFinished {
		return errors.New("message queue is closed")
	}
	md.MessageQueue <- *NewNotificationJob(md.Client, message.Message)
	//fmt.Printf("Dispatching message: %s \n", message.ID)
	return nil
}
