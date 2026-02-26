package grpc

import (
	"context"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/subscription"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/tracing"
	gRPCInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc"
)

type SubscriptionHandler struct {
	pb.UnimplementedSubscriptionServiceServer
	config              *config.Config
	logger              *logging.Logger
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(
	config *config.Config,
	logger *logging.Logger,
	subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		config:              config,
		logger:              logger,
		subscriptionService: subscriptionService,
	}
}

func (h *SubscriptionHandler) GetBookCategorySubscribedUserEmails(
	ctx context.Context,
	req *pb.GetBookCategorySubscribedUserEmailsRequest,
) (*pb.GetBookCategorySubscribedUserEmailsResponse, error) {
	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.GetBookCategorySubscribedUserEmails")
	defer span.End()

	emails, err := h.subscriptionService.GetBookCategorySubscribedUserEmails(ctx, req.GetBookCategory())
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Get book category subscribed user email error", logging.Error(err))
		return nil, gRPCInfrastructure.NewError(err)
	}

	return &pb.GetBookCategorySubscribedUserEmailsResponse{Emails: emails}, nil
}
