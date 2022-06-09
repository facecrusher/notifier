package processor

import (
	"log"
	"notifier/message/domain"
	"notifier/rest"
)

type Job interface {
	Process() error
}

type NotificationJob struct {
	Client  rest.NotifierRestClient
	Message domain.Message
}

func NewNotificationJob(client rest.NotifierRestClient, message string) *NotificationJob {
	return &NotificationJob{
		Client:  client,
		Message: *domain.NewMessage(message),
	}
}

func (j *NotificationJob) Process() error {
	var decode map[string]string
	log.Printf("Sending notification: [id = %s][message = %s]\n", j.Message.ID, j.Message.Message)
	err := j.Client.Post(j.Message, decode)
	if err != nil {
		log.Printf("Error sending notification: [id = %s][message = %s][error = %s]\n", j.Message.ID, j.Message.Message, err.Error())
	}
	return err
}
