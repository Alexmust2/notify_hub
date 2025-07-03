package app

import (
	"fmt"

	"github.com/notify-hub/internal/logger"
	"github.com/notify-hub/internal/notifier"
	"github.com/notify-hub/internal/queue"
)

type ChannelNotification struct {
	Channel        string
	IntegrationKey string
	Receivers      []string
}

type ChannelResult struct {
	Channel      string
	Success      bool
	ErrorMessage string
}

type NotificationUseCase struct {
	notifiers map[string]notifier.Notifier
	queue     queue.Queue
	logger    logger.Logger
}

func NewNotificationUseCase(
	notifiers map[string]notifier.Notifier,
	queue queue.Queue,
	logger logger.Logger,
) *NotificationUseCase {
	return &NotificationUseCase{
		notifiers: notifiers,
		queue:     queue,
		logger:    logger,
	}
}

func (uc *NotificationUseCase) SendNotificationMulti(notifications []ChannelNotification, message string) []ChannelResult {
	var results []ChannelResult

	for _, n := range notifications {
		notifier, exists := uc.notifiers[n.Channel]
		if !exists {
			errMsg := fmt.Sprintf("channel %s not supported", n.Channel)
			uc.logger.Error(errMsg)
			results = append(results, ChannelResult{
				Channel:      n.Channel,
				Success:      false,
				ErrorMessage: errMsg,
			})
			continue
		}

		err := notifier.Send(n.IntegrationKey, n.Receivers, message)
		if err != nil {
			uc.logger.Error(fmt.Sprintf("Failed to send notification to channel %s: %v", n.Channel, err))
			results = append(results, ChannelResult{
				Channel:      n.Channel,
				Success:      false,
				ErrorMessage: err.Error(),
			})
			continue
		}

		uc.logger.Info(fmt.Sprintf("Notification sent to channel %s", n.Channel))
		results = append(results, ChannelResult{
			Channel:      n.Channel,
			Success:      true,
			ErrorMessage: "",
		})
	}

	return results
}
