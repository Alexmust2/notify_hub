package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/notify-hub/internal/app"
	"github.com/notify-hub/internal/config"
	"github.com/notify-hub/internal/grpc"
	"github.com/notify-hub/internal/logger"
	"github.com/notify-hub/internal/notifier"
	"github.com/notify-hub/internal/queue"
	pb "github.com/notify-hub/pkg/proto"
	grpcServer "google.golang.org/grpc"
)

func main() {
	// Инициализация логгера
	logger := logger.New()

	// Загрузка конфигурации
	cfg, err := config.Load("configs/integrations.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация очереди
	queue := queue.NewInMemoryQueue(logger)

	// Инициализация нотификаторов
	telegramNotifier := notifier.NewTelegramNotifier(cfg.Telegram, logger)
	emailNotifier := notifier.NewEmailNotifier(cfg.Email, logger)

	notifiers := map[string]notifier.Notifier{
		"telegram": telegramNotifier,
		"email":    emailNotifier,
	}

	// Инициализация UseCase
	notificationUseCase := app.NewNotificationUseCase(notifiers, queue, logger)

	// Инициализация gRPC сервера
	grpcHandler := grpc.NewNotificationHandler(notificationUseCase, logger)

	server := grpcServer.NewServer()
	pb.RegisterNotificationServiceServer(server, grpcHandler)

	// Запуск сервера
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Запуск worker'а очереди
	ctx, cancel := context.WithCancel(context.Background())
	go queue.StartWorker(ctx)

	logger.Info("Notify Hub started on :8080")

	// Graceful shutdown
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
	cancel()
	server.GracefulStop()
}
