package processor

import (
	"sync"

	"github.com/facecrusher/notifier/message/domain"
	"github.com/facecrusher/notifier/rest/client"
)

type Sender interface {
	Send(message domain.Message) error
}

type MessageSender struct {
	RestClient       client.RestClient
	done             sync.WaitGroup
	readyPool        chan chan NotificationJob
	assignedJobQueue chan NotificationJob
	quit             chan bool
}

func NewMessageSender(restClient client.RestClient, readyPool chan chan NotificationJob, done sync.WaitGroup) *MessageSender {
	return &MessageSender{
		RestClient:       restClient,
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan NotificationJob),
		quit:             make(chan bool),
	}
}

func (ms *MessageSender) Start() {
	ms.done.Add(1)
	go func() {
		for {
			ms.readyPool <- ms.assignedJobQueue
			select {
			case notificationJob := <-ms.assignedJobQueue:
				notificationJob.Process()
			case <-ms.quit:
				ms.done.Done()
				return
			}
		}
	}()
}

func (ms *MessageSender) Stop() {
	ms.quit <- true
}
