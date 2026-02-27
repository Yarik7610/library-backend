package subscription

import (
	"context"
	"time"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/subscription"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Yarik7610/library-backend-common/microservice"
)

type Client interface {
	GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error)
}

type client struct {
	gRPCClient pb.SubscriptionServiceClient
}

func NewClient() (Client, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		microservice.SUBSCRIPTIONS_GRPC_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}

	return &client{gRPCClient: pb.NewSubscriptionServiceClient(conn)}, conn, nil
}

func (c *client) GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, err := c.gRPCClient.GetBookCategorySubscribedUserEmails(ctx, &pb.GetBookCategorySubscribedUserEmailsRequest{BookCategory: bookCategory})
	if err != nil {
		return nil, err
	}
	return resp.GetEmails(), nil
}
