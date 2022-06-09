package processor

import (
	"notifier/rest"
	"sync"
)

const (
	DEFAULT_SENDERS    = 5
	DEFAULT_QUEUE_SIZE = 10
)

type MessageQueue struct {
	Options        Options
	internalQueue  chan NotificationJob
	readyPool      chan chan NotificationJob
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
	readyPool := make(chan chan NotificationJob, options.MaxSenders)   // set bidirectional channel for available senders
	return &MessageQueue{
		Options:        *options,
		internalQueue:  make(chan NotificationJob),
		readyPool:      readyPool,
		senders:        createSenders(options.MaxSenders, client, readyPool, sendersStopped),
		sendersStopped: sendersStopped,
		queueStopped:   sync.WaitGroup{},
		quit:           make(chan bool),
	}
}

func (mq *MessageQueue) GetMessageQueue() chan NotificationJob {
	return mq.internalQueue
}

func (mq *MessageQueue) Start() {
	mq.startSenders()
	go mq.initDispatch()
}

func (mq *MessageQueue) initDispatch() {
	mq.queueStopped.Add(1)
	for {
		select {
		case notificationJob := <-mq.internalQueue: // a notification job is present in the queue
			senderChannel := <-mq.readyPool  // look for an available sender
			senderChannel <- notificationJob // send the job to the senders channel for processing
		case <-mq.quit:
			mq.stopSenders() //stop senders and wait for them to be done with any on going message delivery
			mq.sendersStopped.Wait()
			mq.queueStopped.Done() // stop the message queue
			return
		}
	}
}

func (mq *MessageQueue) Stop() {
	if len(mq.internalQueue) == 0 {
		mq.quit <- true
		mq.queueStopped.Wait()
	}
	mq.queueStopped.Wait()
}

func createSenders(senderAmount int, client *rest.NotifierRestClient,
	readyPool chan chan NotificationJob, done sync.WaitGroup) []*MessageSender {
	var senders []*MessageSender
	for i := 0; i < senderAmount; i++ {
		senders = append(senders, NewMessageSender(*client, readyPool, done))
	}
	return senders
}

func (mq *MessageQueue) startSenders() {
	for i := 0; i < len(mq.senders); i++ {
		mq.senders[i].Start()
	}
}

func (mq *MessageQueue) stopSenders() {
	for i := 0; i < len(mq.senders); i++ {
		mq.senders[i].Stop()
	}
}
