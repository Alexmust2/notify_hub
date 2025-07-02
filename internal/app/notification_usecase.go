package app

import (
	"fmt"

	"github.com/notify-hub/internal/logger"
	"github.com/notify-hub/internal/notifier"
	"github.com/notify-hub/internal/queue"
)

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

func (uc *NotificationUseCase) SendNotification(channel, integrationKey, receiver, message string) error {
	notifier, exists := uc.notifiers[channel]
	if !exists {
		return fmt.Errorf("channel %s not supported", channel)
	}

	// Для синхронной отправки (можно изменить на асинхронную через очередь)
	err := notifier.Send(integrationKey, receiver, message)
	if err != nil {
		uc.logger.Error(fmt.Sprintf("Failed to send notification: %v", err))
		return err
	}

	return nil
}

func (uc *NotificationUseCase) SendNotificationAsync(channel, integrationKey, receiver, message string) {
	task := queue.NotificationTask{
		Channel:        channel,
		IntegrationKey: integrationKey,
		Receiver:       receiver,
		Message:        message,
		Callback: func(err error) {
			if err != nil {
				uc.logger.Error(fmt.Sprintf("Async notification failed: %v", err))
			} else {
				uc.logger.Info("Async notification sent successfully")
			}
		},
	}

	uc.queue.Enqueue(task)
}