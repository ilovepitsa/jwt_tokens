package notification

import "log"

type emailNotificator struct {
}

func NewEmailNotificator() *emailNotificator {
	return &emailNotificator{}
}

func (en *emailNotificator) Notificate(text string) {
	log.Printf("\n\n%s\n\n", text)
}
