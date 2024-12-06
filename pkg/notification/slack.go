package notification

import "github.com/pchchv/goddns/internal/settings"

type SlackNotification struct {
	conf *settings.Settings
}

func NewSlackNotification(conf *settings.Settings) INotification {
	return &SlackNotification{conf: conf}
}