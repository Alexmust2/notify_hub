package notifier

import (
	"fmt"
	"net/smtp"

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

func (e *EmailNotifier) Send(integrationKey, to, message string) error {
	cfg, exists := e.configs[integrationKey]
	if !exists {
		return fmt.Errorf("email integration key %s not found", integrationKey)
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	
	msg := fmt.Sprintf("To: %s\r\nSubject: Notification\r\n\r\n%s", to, message)
	
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	
	err := smtp.SendMail(addr, auth, cfg.Username, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.Info(fmt.Sprintf("Email sent successfully to %s via %s", to, integrationKey))
	return nil
}