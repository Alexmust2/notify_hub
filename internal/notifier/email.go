package notifier

import (
	"fmt"

	mail "github.com/go-mail/mail/v2"
	"github.com/notify-hub/internal/config"
	"github.com/notify-hub/internal/logger"
)

type EmailNotifier struct {
	configs map[string]config.EmailConfig
	logger  logger.Logger
}

func NewEmailNotifier(configs map[string]config.EmailConfig, logger logger.Logger) *EmailNotifier {
	return &EmailNotifier{
		configs: configs,
		logger:  logger,
	}
}

func (e *EmailNotifier) Send(integrationKey string, to []string, message string, metadata map[string]string) error {
	cfg, exists := e.configs[integrationKey]
	if !exists {
		return fmt.Errorf("email integration key %s not found", integrationKey)
	}

	subject := "Notification"
	if customSubject, ok := metadata["subject"]; ok {
		subject = customSubject
	}

	m := mail.NewMessage()
	m.SetHeader("From", cfg.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := mail.NewDialer(cfg.Host, 465, cfg.Username, cfg.Password)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Info(fmt.Sprintf("Email sent successfully to %v via %s", to, integrationKey))
	return nil
}
