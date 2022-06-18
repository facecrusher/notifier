package processor

import (
	"fmt"
	"log"
	"time"

	"github.com/facecrusher/notifier/message/domain"

	"github.com/facecrusher/notifier/rest/client"
)

const (
	RETRY_INTERVAL = 2 * client.DEFAULT_TIMEOUT
)

//go:generate mockgen -source=./dispatcher.go -destination=./mock/dispatcher.go -package=mock
type Dispatcher interface {
	Dispatch(message domain.Message) error
}

type MessageDispatcher struct {
	client       client.RestClient
	messageQueue chan NotificationJob
	interval     time.Duration
}

func NewDispatcher(mq MessageQueue, c client.RestClient, interval *time.Duration) Dispatcher {
	return &MessageDispatcher{
		client:       c,
		messageQueue: mq.GetMessageQueue(),
		interval:     getInterval(interval),
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	messageJob := *NewNotificationJob(md.client, message.Message, md.interval)
	select {
	case md.messageQueue <- messageJob:
	default:
		log.Printf("Processing queue full. Putting message [%s] in retry queue.", messageJob.Message.Message)
		return md.Retry(messageJob)
	}
	return nil
}

// Retry attempts to re-queue a message in case the buffer was full at the time message was dispatched
func (md *MessageDispatcher) Retry(job NotificationJob) error {
	time.Sleep(RETRY_INTERVAL)
	select {
	case md.messageQueue <- job:
	default:
		return fmt.Errorf("processing queue full when retried. Message [%s] discarded", job.Message.Message)
	}
	return nil
}

func getInterval(inputInterval *time.Duration) (interval time.Duration) {
	if inputInterval == nil {
		return 0
	}
	return *inputInterval
}
