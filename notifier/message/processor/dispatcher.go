package processor

import (
	"errors"
	"notifier/notifier/message/domain"
)

type Dispatcher interface {
	Dispatch(message domain.Message) error
}

type MessageDispatcher struct {
	MessageQueue chan domain.Message
	IsFinished   bool
}

func NewDispatcher(mq MessageQueue) *MessageDispatcher {
	return &MessageDispatcher{
		MessageQueue: mq.GetMessageQueue(),
		IsFinished:   false,
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	if md.IsFinished {
		return errors.New("message queue is closed")
	}
	md.MessageQueue <- message
	return nil
}
