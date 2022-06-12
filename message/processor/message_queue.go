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
	availablePool  chan chan NotificationJob
	senders        []*MessageSender
	sendersStopped sync.WaitGroup
	QueueStopped   sync.WaitGroup
	quit           chan bool
}

type Options struct {
	MaxSenders   int
	MaxQueueSize int
}

func NewMessageQueue(url string, options *Options, client rest.NotifierRestClient) *MessageQueue {
	// Set default queue options if none are provided
	if options == nil {
		options = getDefaultOptions()
	}
	sendersStopped := sync.WaitGroup{}                                   // set wait group for stopped senders
	availablePool := make(chan chan NotificationJob, options.MaxSenders) // set pool to maintain available senders
	return &MessageQueue{
		Options:        *options,
		internalQueue:  make(chan NotificationJob, options.MaxQueueSize),
		availablePool:  availablePool,
		senders:        createSenders(options.MaxSenders, client, availablePool, sendersStopped),
		sendersStopped: sendersStopped,
		QueueStopped:   sync.WaitGroup{},
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

// Stop shuts down the message queue, checking that all messages in the buffer are processed and that
// before sending the quit signal.
func (mq *MessageQueue) Stop() {
	mq.quit <- true
	mq.QueueStopped.Wait()
}

func (mq *MessageQueue) initDispatch() {
	mq.QueueStopped.Add(1)
	for {
		select {
		case notificationJob := <-mq.internalQueue: // a notification job is present in the queue
			senderChannel := <-mq.availablePool // look for an available sender in the availablePool
			senderChannel <- notificationJob    // send the job to the selected sender assignedJobQueue for processing
		case <-mq.quit:
			mq.processPendingMessages() //process any pending notification jobs in the queue buffer
			mq.stopSenders()            // stop senders
			mq.sendersStopped.Wait()    // wait for senders to finish any on going work
			mq.QueueStopped.Done()      // stop the message queue
			return
		}
	}
}

func (mq *MessageQueue) processPendingMessages() {
	for {
		select {
		case notificationJob := <-mq.internalQueue:
			senderChannel := <-mq.availablePool
			senderChannel <- notificationJob
		default:
			return
		}
	}
}

func getDefaultOptions() *Options {
	return &Options{MaxSenders: DEFAULT_SENDERS, MaxQueueSize: DEFAULT_QUEUE_SIZE}
}

func createSenders(senderAmount int, client rest.NotifierRestClient,
	readyPool chan chan NotificationJob, done sync.WaitGroup) []*MessageSender {
	var senders []*MessageSender
	for i := 0; i < senderAmount; i++ {
		senders = append(senders, NewMessageSender(client, readyPool, done))
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
