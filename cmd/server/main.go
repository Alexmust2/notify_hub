package main

import (
	"log"
	"net"

	"notify_hub/internal/config"
	"notify_hub/internal/delivery/grpc"
	"notify_hub/internal/service"
	"notify_hub/internal/telegram"
	pb "notify_hub/proto"

	"github.com/joho/godotenv"

	grpcGo "google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	cfg := config.Load()

	tgClient := telegram.NewClient(cfg.TelegramBotToken)
	notifService := service.NewNotificationService(tgClient)

	grpcServer := grpc.NewServer(notifService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpcGo.NewServer() // Это grpc из Google
	pb.RegisterNotificationServiceServer(s, grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
