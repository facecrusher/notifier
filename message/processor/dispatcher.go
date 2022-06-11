package processor

import (
	"fmt"
	"log"
	"notifier/message/domain"
	"notifier/rest"
	"time"
)

const (
	RETRY_INTERVAL = 2 * rest.DEFAULT_TIMEOUT
)

type Dispatcher interface {
	Dispatch(message domain.Message) error
}

type MessageDispatcher struct {
	Client       rest.NotifierRestClient
	MessageQueue chan NotificationJob
	Interval     time.Duration
}

func NewDispatcher(mq MessageQueue, c rest.NotifierRestClient, interval *time.Duration) *MessageDispatcher {
	return &MessageDispatcher{
		Client:       c,
		MessageQueue: mq.GetMessageQueue(),
		Interval:     getInterval(interval),
	}
}

func (md *MessageDispatcher) Dispatch(message domain.Message) error {
	messageJob := *NewNotificationJob(md.Client, message.Message, md.Interval)
	select {
	case md.MessageQueue <- messageJob:
	default:
		log.Printf("Processing queue full. Putting message [%s] in retry queue.", messageJob.Message.Message)
		return md.Retry(messageJob)
	}
	return nil
}

// Retry attempts to re-queue a message in case the buffer was full full at the time message was received
func (md *MessageDispatcher) Retry(job NotificationJob) error {
	time.Sleep(RETRY_INTERVAL)
	select {
	case md.MessageQueue <- job:
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
