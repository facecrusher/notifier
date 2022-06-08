package processor

import (
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
	"sync"
)

const (
	DEFAULT_SENDERS    = 3
	DEFAULT_QUEUE_SIZE = 9
)

type MessageQueue struct {
	Options        Options
	internalQueue  chan domain.Message
	readyPool      chan chan domain.Message
	senders        []*MessageSender
	sendersStopped sync.WaitGroup
	queueStopped   sync.WaitGroup
	quit           chan bool
}

type Options struct {
	MaxSenders   int
	MaxQueueSize int
}

func NewMessageQueue(url string, options *Options) *MessageQueue {
	// Set default queue options if none are provided
	if options == nil {
		options = &Options{MaxSenders: DEFAULT_SENDERS, MaxQueueSize: DEFAULT_QUEUE_SIZE}
	}
	client := rest.NewNotifierRestClient(url, make(map[string]string)) // set rest client
	sendersStopped := sync.WaitGroup{}                                 // set wait group for stopped senders
	readyPool := make(chan chan domain.Message, options.MaxSenders)    // set wait group for available pool of senders
	return &MessageQueue{
		Options:        *options,
		internalQueue:  make(chan domain.Message, options.MaxQueueSize),
		readyPool:      readyPool,
		senders:        createSenders(options.MaxSenders, client, readyPool, sendersStopped),
		sendersStopped: sendersStopped,
		queueStopped:   sync.WaitGroup{},
		quit:           make(chan bool),
	}
}

func (mq *MessageQueue) GetMessageQueue() chan domain.Message {
	return mq.internalQueue
}

func (mq *MessageQueue) Start() {
	for i := 0; i < len(mq.senders); i++ {
		mq.senders[i].Start()
	}
	mq.queueStopped.Add(1)
	go func() {
		for {
			select {
			case message := <-mq.internalQueue: // a message is present in the queue
				//fmt.Printf("Processing message: %s \n", message.ID)
				senderChannel := <-mq.readyPool // look for an available sender
				senderChannel <- message        // send the message to the senders channel for processing
			case <-mq.quit:
				mq.stopSenders()       //stop senders and wait for them to be done with any on going message delivery
				mq.queueStopped.Done() // stop the message queue
				return
			}
		}
	}()
}

func (mq *MessageQueue) Stop() {
	mq.quit <- true
	mq.queueStopped.Wait()
}

func createSenders(senderAmount int, client *rest.NotifierRestClient,
	readyPool chan chan domain.Message, done sync.WaitGroup) []*MessageSender {
	var senders []*MessageSender
	for i := 0; i < senderAmount; i++ {
		senders = append(senders, NewMessageSender(*client, readyPool, done))
	}
	return senders
}

func (mq *MessageQueue) stopSenders() {
	for _, sender := range mq.senders {
		sender.Stop()
	}
	mq.sendersStopped.Wait()
}
