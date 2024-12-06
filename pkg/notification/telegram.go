package notification

import "github.com/pchchv/goddns/internal/settings"

type TelegramNotification struct {
	conf *settings.Settings
}

func NewTelegramNotification(conf *settings.Settings) INotification {
	return &TelegramNotification{conf: conf}
}
