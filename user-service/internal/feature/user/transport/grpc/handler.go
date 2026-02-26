package grpc

import (
	"context"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/user"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/tracing"
	grpcInfrastructure "github.com/Yarik7610/library-backend/user-service/internal/infrastructure/transport/grpc"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	config      *config.Config
	logger      *logging.Logger
	userService service.UserService
}

func NewUserHandler(
	config *config.Config,
	logger *logging.Logger,
	userService service.UserService) *UserHandler {
	return &UserHandler{
		config:      config,
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) GetEmailsByUserIDs(ctx context.Context, req *pb.GetEmailsByUserIDsRequest) (*pb.GetEmailsByUserIDsResponse, error) {
	userIDs := make([]uint, len(req.GetUserIds()))
	for i, userID := range req.UserIds {
		userIDs[i] = uint(userID)
	}

	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.GetEmailsByUserIDs")
	defer span.End()

	emails, err := h.userService.GetEmailsByUserIDs(ctx, userIDs)
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Get emails by user IDs error", logging.Error(err))
		return nil, grpcInfrastructure.NewError(err)
	}

	return &pb.GetEmailsByUserIDsResponse{Emails: emails}, nil
}
