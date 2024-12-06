package notification

type INotification interface {
	Send(domain, currentIP string) error
}
