package grpc

import (
	"context"
	"notify_hub/internal/service"
	pb "notify_hub/proto"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer
	service *service.NotificationService
}

func NewServer(service *service.NotificationService) *Server {
	return &Server{service: service}
}

func (s *Server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	var err error
	switch req.MessageType {
	case pb.MessageType_TEXT:
		err = s.service.SendToTelegram(req.TelegramUserId, req.Content)
	case pb.MessageType_PHOTO:
		err = s.service.SendPhotoToTelegram(req.TelegramUserId, req.Content, req.FileData, req.Filename)
	case pb.MessageType_DOCUMENT:
		err = s.service.SendDocumentToTelegram(req.TelegramUserId, req.Content, req.FileData, req.Filename)
	case pb.MessageType_AUDIO:
		err = s.service.SendAudioToTelegram(req.TelegramUserId, req.Content, req.FileData, req.Filename)
	default:
		return &pb.SendNotificationResponse{Success: false}, nil
	}

	if err != nil {
		return &pb.SendNotificationResponse{Success: false}, err
	}
	return &pb.SendNotificationResponse{Success: true}, nil
}
