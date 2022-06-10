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
}

func NewDispatcher(mq MessageQueue, c rest.NotifierRestClient, interval time.Duration) *MessageDispatcher {
	return &MessageDispatcher{
		Client:       c,
		MessageQueue: mq.GetMessageQueue(),
		Interval:     interval,
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	if isQueueClosed(md.MessageQueue) {
		return errors.New("message queue is closed")
	}
	md.MessageQueue <- *NewNotificationJob(md.Client, message.Message, md.Interval)
	return nil
}

func isQueueClosed(messageQueue chan NotificationJob) bool {
	select {
	case <-messageQueue:
		return true
	default:
		return false
	}
}
