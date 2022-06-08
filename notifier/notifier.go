package notifier

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
