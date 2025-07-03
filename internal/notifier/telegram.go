package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/notify-hub/internal/config"
	"github.com/notify-hub/internal/logger"
)

type TelegramNotifier struct {
	configs map[string]config.TelegramConfig
	logger  logger.Logger
}

func NewTelegramNotifier(configs map[string]config.TelegramConfig, logger logger.Logger) *TelegramNotifier {
	return &TelegramNotifier{
		configs: configs,
		logger:  logger,
	}
}

type telegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (t *TelegramNotifier) Send(integrationKey string, to []string, message string) error {
	cfg, exists := t.configs[integrationKey]
	if !exists {
		return fmt.Errorf("telegram integration key %s not found", integrationKey)
	}

	for _, receiver := range to {
		msg := telegramMessage{
			ChatID: receiver,
			Text:   message,
		}

		jsonData, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal telegram message: %w", err)
		}

		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.Token)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}
		defer resp.Body.Close()

		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return fmt.Errorf("failed to read telegram response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, responseBody.String())
		}

		var telegramResp struct {
			Ok          bool   `json:"ok"`
			Description string `json:"description"`
		}

		if err := json.Unmarshal(responseBody.Bytes(), &telegramResp); err != nil {
			return fmt.Errorf("failed to decode telegram response: %w. Raw response: %s", err, responseBody.String())
		}

		if !telegramResp.Ok {
			return fmt.Errorf("telegram API error: %s. Raw response: %s", telegramResp.Description, responseBody.String())
		}

		t.logger.Info(fmt.Sprintf("Telegram message sent successfully to %s via %s", receiver, integrationKey))
	}

	return nil
}
