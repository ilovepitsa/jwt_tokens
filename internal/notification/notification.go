package notification

type Notifier interface {
	Notificate(text string)
}
