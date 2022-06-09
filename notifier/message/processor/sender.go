package processor

import (
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
	"sync"
)

type Sender interface {
	Send(message domain.Message) error
}

type MessageSender struct {
	RestClient       rest.NotifierRestClient
	done             sync.WaitGroup
	readyPool        chan chan NotificationJob
	assignedJobQueue chan NotificationJob
	quit             chan bool
}

func NewMessageSender(restClient rest.NotifierRestClient, readyPool chan chan NotificationJob, done sync.WaitGroup) *MessageSender {
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
