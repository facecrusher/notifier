package notifier

import (
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

func NewNotifier(url string, interval *time.Duration, options *processor.Options) *Notifier {
	client := rest.NewNotifierRestClient(url, make(map[string]string))
	queue := *processor.NewMessageQueue(url, options)
	dispatcher := processor.NewDispatcher(queue, *client, interval)
	return &Notifier{
		url:               url,
		messageQueue:      queue,
		interval:          *interval,
		messageDispatcher: dispatcher,
	}
}
