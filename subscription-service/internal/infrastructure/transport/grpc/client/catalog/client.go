package catalog

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/catalog"
	grpcInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	BookCategoryExists(ctx context.Context, bookCategory string) (bool, error)
}

type client struct {
	gRPCClient pb.CatalogServiceClient
}

func NewClient() (Client, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		microservice.CATALOG_GRPC_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, nil, err
	}

	return &client{gRPCClient: pb.NewCatalogServiceClient(conn)}, conn, nil
}

func (c *client) BookCategoryExists(ctx context.Context, bookCategory string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, err := c.gRPCClient.BookCategoryExists(ctx, &pb.BookCategoryExistsRequest{BookCategory: bookCategory})
	if err != nil {
		return false, grpcInfrastructure.ToInfrastctureError(err)
	}
	return resp.GetExists(), nil
}
