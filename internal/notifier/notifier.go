package notifier

type Notifier interface {
	Send(integrationKey string, to []string, message string, metadata map[string]string) error
}
