package notification

type INotification interface {
	Send(domain, currentIP string) error
}

type INotificationManager interface {
	Send(string, string)
}
