package grpc

import (
	"context"
	"fmt"

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
	h.logger.Info("Received multi-channel notification request")

	var notifications []app.ChannelNotification
	for _, n := range req.Notifications {
		notifications = append(notifications, app.ChannelNotification{
			Channel:        n.Channel,
			IntegrationKey: n.IntegrationKey,
			Receivers:      n.Receivers,
		})
	}

	results := h.useCase.SendNotificationMulti(notifications, req.Message)

	for _, r := range results {
		h.logger.Info(fmt.Sprintf("Channel %s success: %v, error: %s", r.Channel, r.Success, r.ErrorMessage))
	}

	var protoResults []*pb.ChannelResult
	for _, r := range results {
		protoResults = append(protoResults, &pb.ChannelResult{
			Channel:      r.Channel,
			Success:      r.Success,
			ErrorMessage: r.ErrorMessage,
		})
	}

	return &pb.SendNotificationResponse{
		Results: protoResults,
	}, nil
}
