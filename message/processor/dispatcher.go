package processor

import (
	"errors"
	"notifier/message/domain"
	"notifier/rest"
	"time"
)

type Dispatcher interface {
	Dispatch(message domain.Message) error
}

type MessageDispatcher struct {
	Client       rest.NotifierRestClient
	MessageQueue chan NotificationJob
	Interval     time.Duration
	IsFinished   bool
}

func NewDispatcher(mq MessageQueue, c rest.NotifierRestClient, interval time.Duration) *MessageDispatcher {
	return &MessageDispatcher{
		Client:       c,
		MessageQueue: mq.GetMessageQueue(),
		Interval:     interval,
		IsFinished:   false,
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	if md.IsFinished {
		return errors.New("message queue is closed")
	}
	md.MessageQueue <- *NewNotificationJob(md.Client, message.Message, md.Interval)
	return nil
}
