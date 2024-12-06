package notification

import "github.com/pchchv/goddns/internal/settings"

type DiscordNotification struct {
	conf *settings.Settings
}

func NewDiscordNotification(conf *settings.Settings) INotification {
	return &DiscordNotification{conf: conf}
}
