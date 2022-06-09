package processor

import (
	"fmt"
	"notifier/notifier/message/domain"
	"notifier/notifier/rest"
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
	fmt.Printf("Sending notification: [id = %s][message = %s]\n", j.Message.ID, j.Message.Message)
	return j.Client.Post(j.Message, decode)
}
