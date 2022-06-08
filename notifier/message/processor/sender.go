package processor

import (
	"fmt"
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
	readyPool        chan chan domain.Message
	assignedJobQueue chan domain.Message
	quit             chan bool
}

func NewMessageSender(restClient rest.NotifierRestClient, readyPool chan chan domain.Message, done sync.WaitGroup) *MessageSender {
	return &MessageSender{
		RestClient:       restClient,
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan domain.Message),
		quit:             make(chan bool),
	}
}

func (ms *MessageSender) Send(message domain.Message) error {
	var decode map[string]string
	fmt.Printf("Sending message: [id = %s] [message = %s] \n", message.ID, message.Message)
	return ms.RestClient.Post(message, decode)
}

func (ms *MessageSender) Start() {
	ms.done.Add(1)
	go func() {
		for {
			ms.readyPool <- ms.assignedJobQueue
			select {
			case message := <-ms.assignedJobQueue:
				ms.Send(message)
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
