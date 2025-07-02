package queue

import (
	"context"
	"time"

	"fmt"

	"github.com/notify-hub/internal/logger"
)

type NotificationTask struct {
	Channel        string
	IntegrationKey string
	Receiver       string
	Message        string
	Callback       func(error)
}

type Queue interface {
	Enqueue(task NotificationTask)
	StartWorker(ctx context.Context)
}

type InMemoryQueue struct {
	tasks  chan NotificationTask
	logger logger.Logger
}

func NewInMemoryQueue(logger logger.Logger) *InMemoryQueue {
	return &InMemoryQueue{
		tasks:  make(chan NotificationTask, 1000),
		logger: logger,
	}
}

func (q *InMemoryQueue) Enqueue(task NotificationTask) {
	select {
	case q.tasks <- task:
		q.logger.Info("Task enqueued successfully")
	default:
		q.logger.Error("Queue is full, dropping task")
		if task.Callback != nil {
			task.Callback(fmt.Errorf("queue is full"))
		}
	}
}

func (q *InMemoryQueue) StartWorker(ctx context.Context) {
	q.logger.Info("Queue worker started")

	for {
		select {
		case <-ctx.Done():
			q.logger.Info("Queue worker stopped")
			return
		case <-q.tasks:
			// Здесь должна быть логика обработки задачи
			// В реальном приложении это будет передано в NotificationUseCase
			q.logger.Info("Processing task...")
			time.Sleep(100 * time.Millisecond) // Имитация обработки
		}
	}
}
