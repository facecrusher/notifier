package notifier

import (
	"log"
	"time"

	"github.com/facecrusher/notifier/message/processor"
	"github.com/facecrusher/notifier/rest/client"

	"github.com/facecrusher/notifier/message/domain"
)

type Notifier struct {
	url               string
	messageQueue      processor.MessageQueue
	messageDispatcher processor.Dispatcher
	interval          time.Duration
}

func NewNotifier(url string, interval *time.Duration, options *processor.Options, headers *map[string]string) *Notifier {
	client := client.NewNotifierRestClient(url, headers)
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

func (n *Notifier) Start() {
	n.messageQueue.Start()
}

func (n *Notifier) Stop() {
	log.Printf("Sending shut down signal to all workers")
	n.messageQueue.Stop()
}
