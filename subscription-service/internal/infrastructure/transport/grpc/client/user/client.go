package user

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/user"
	grpcInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error)
}

type client struct {
	gRPCClient pb.UserServiceClient
}

func NewClient() (Client, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		microservice.USER_GRPC_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}

	return &client{gRPCClient: pb.NewUserServiceClient(conn)}, conn, nil
}

func (c *client) GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error) {
	userIDsUint64 := make([]uint64, len(userIDs))
	for i, userID := range userIDs {
		userIDsUint64[i] = uint64(userID)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, err := c.gRPCClient.GetEmailsByUserIDs(ctx, &pb.GetEmailsByUserIDsRequest{UserIds: userIDsUint64})
	if err != nil {
		return nil, grpcInfrastructure.ToInfrastctureError(err)
	}
	return resp.GetEmails(), nil
}
