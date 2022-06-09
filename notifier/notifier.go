package notifier

import (
	"fmt"
	"notifier/message/domain"
	"notifier/message/processor"
	"notifier/rest"
)

type Notifier struct {
	URL string
}

func NewNotifier(URL string) *Notifier {
	return &Notifier{
		URL: URL,
	}
}

func SendMessage(url, message string) {

}

func main() {
	url := "https://eouss1txxyn5t7x.m.pipedream.net"
	client := rest.NewNotifierRestClient(url, make(map[string]string))
	messageQueue := processor.NewMessageQueue(url, nil)
	dispatcher := processor.NewDispatcher(*messageQueue, *client)

	messageQueue.Start()
	defer messageQueue.Stop()

	for i := 0; i < 20; i++ {
		message := domain.NewMessage(fmt.Sprintf("This is message number %d", i))
		dispatcher.Dispatch(*message)
	}
}
