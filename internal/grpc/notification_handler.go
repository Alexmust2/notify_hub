package grpc

import (
	"context"

	"github.com/notify-hub/internal/app"
	"github.com/notify-hub/internal/logger"
	pb "github.com/notify-hub/pkg/proto"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	useCase *app.NotificationUseCase
	logger  logger.Logger
}

func NewNotificationHandler(useCase *app.NotificationUseCase, logger logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *NotificationHandler) SendNotification(
	ctx context.Context,
	req *pb.SendNotificationRequest,
) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received notification request")

	err := h.useCase.SendNotification(
		req.Channel,
		req.IntegrationKey,
		req.Receiver,
		req.Message,
	)

	if err != nil {
		return &pb.SendNotificationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.SendNotificationResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}