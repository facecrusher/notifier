package processor

import (
	"log"
	"time"

	"github.com/facecrusher/notifier/message/domain"
	"github.com/facecrusher/notifier/rest/client"
)

type Job interface {
	Process() error
}

type NotificationJob struct {
	Client    client.RestClient
	Message   domain.Message
	Interval  time.Duration
	processed *chan string //optional chan for testing
}

func NewNotificationJob(client client.RestClient, message string, interval time.Duration, processedChan *chan string) *NotificationJob {
	return &NotificationJob{
		Client:    client,
		Message:   *domain.NewMessage(message),
		Interval:  interval,
		processed: processedChan,
	}
}

func (j *NotificationJob) Process() error {
	log.Printf("Sending notification: [id = %s][message = %s]\n", j.Message.ID, j.Message.Message)
	var decode map[string]string
	time.Sleep(j.Interval)
	err := j.Client.Post(j.Message, decode)
	if err != nil {
		log.Printf("Error sending notification: [id = %s][message = %s][error = %s]\n", j.Message.ID, j.Message.Message, err.Error())
	}
	j.reportProcessed()
	return err
}

func (j *NotificationJob) reportProcessed() {
	if j.processed != nil {
		*j.processed <- j.Message.Message
	}
}
