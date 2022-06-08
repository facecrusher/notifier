package processor

import (
	"context"
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
	"sync"
)

const (
	DefaultSenders   = 10
	DefaultQueueSize = 1000
)

type MessageQueue struct {
	Options       Options
	internalQueue chan domain.Message
	sender        *MessageSender
}

type Options struct {
	MaxSenders   int
	MaxQueueSize int
}

func NewMessageQueue(url string, options *Options) *MessageQueue {
	headers := make(map[string]string)
	client := rest.NewNotifierRestClient(url, headers)

	if options == nil {
		options = &Options{
			MaxSenders:   DefaultSenders,
			MaxQueueSize: DefaultQueueSize,
		}
	}
	return &MessageQueue{
		Options:       *options,
		internalQueue: make(chan domain.Message, options.MaxQueueSize),
		sender:        NewMessageSender(*client),
	}
}

func (mq *MessageQueue) GetMessageQueue() chan domain.Message {
	return mq.internalQueue
}

func (mq *MessageQueue) Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	for i := 0; i < mq.Options.MaxSenders; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case m := <-mq.internalQueue:
					mq.sender.Send(m)
					return
				}
			}
		}()
	}
}
