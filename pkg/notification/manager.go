package notification

import (
	"log"

	"github.com/pchchv/goddns/internal/settings"
)

const (
	Email    = "email"
	Slack    = "slack"
	Discord  = "discord"
	Telegram = "telegram"
	Pushover = "pushover"
)

type INotification interface {
	Send(domain, currentIP string) error
}

type INotificationManager interface {
	Send(string, string)
}

type notificationManager struct {
	notifications map[string]INotification
}

func (n *notificationManager) Send(domain, currentIP string) {
	for _, sender := range n.notifications {
		if err := sender.Send(domain, currentIP); err != nil {
			log.Fatalf("Send notification with error: %e", err)
		}
	}
}

func initNotifications(conf *settings.Settings) map[string]INotification {
	notificationMap := map[string]INotification{}
	if conf.Notify.Mail.Enabled {
		notificationMap[Email] = NewEmailNotification(conf)
	}

	if conf.Notify.Telegram.Enabled {
		notificationMap[Telegram] = NewTelegramNotification(conf)
	}

	if conf.Notify.Discord.Enabled {
		notificationMap[Discord] = NewDiscordNotification(conf)
	}

	if conf.Notify.Slack.Enabled {
		notificationMap[Slack] = NewSlackNotification(conf)
	}

	if conf.Notify.Pushover.Enabled {
		notificationMap[Pushover] = NewPushoverNotification(conf)
	}

	return notificationMap
}
