package notification

import "github.com/pchchv/goddns/internal/settings"

type PushoverNotification struct {
	conf *settings.Settings
}

func NewPushoverNotification(conf *settings.Settings) INotification {
	return &PushoverNotification{conf: conf}
}
