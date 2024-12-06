package notification

import "log"

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
