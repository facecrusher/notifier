package processor

import (
	"log"
	"notifier/message/domain"
	"notifier/rest"
	"time"
)

type Job interface {
	Process() error
}

type NotificationJob struct {
	Client   rest.NotifierRestClient
	Message  domain.Message
	Interval time.Duration
}

func NewNotificationJob(client rest.NotifierRestClient, message string, interval time.Duration) *NotificationJob {
	return &NotificationJob{
		Client:   client,
		Message:  *domain.NewMessage(message),
		Interval: interval,
	}
}

func (j *NotificationJob) Process() error {
	log.Printf("Sending notification: [id = %s][message = %s]\n", j.Message.ID, j.Message.Message)
	time.Sleep(j.Interval)
	// var decode map[string]string
	// err := j.Client.Post(j.Message, decode)
	// if err != nil {
	// 	log.Printf("Error sending notification: [id = %s][message = %s][error = %s]\n", j.Message.ID, j.Message.Message, err.Error())
	// }
	// return err
	return nil
}
