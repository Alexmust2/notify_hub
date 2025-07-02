package notifier

type Notifier interface {
	Send(integrationKey, to, message string) error
}