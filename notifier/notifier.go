package notifier

import (
	"notifier/message/domain"
	"notifier/message/processor"
	"notifier/rest"
	"time"
)

type Notifier struct {
	url               string
	messageQueue      processor.MessageQueue
	messageDispatcher processor.Dispatcher
	interval          time.Duration
}

func NewNotifier(url string, interval *time.Duration, options *processor.Options, headers *map[string]string) *Notifier {
	client := rest.NewNotifierRestClient(url, headers)
	queue := *processor.NewMessageQueue(url, options, *client)
	dispatcher := processor.NewDispatcher(queue, *client, interval)
	return &Notifier{
		url:               url,
		messageQueue:      queue,
		interval:          *interval,
		messageDispatcher: dispatcher,
	}
}

func (n *Notifier) Notify(messageString string) error {
	message := domain.NewMessage(messageString)
	return n.messageDispatcher.Dispatch(*message)
}
