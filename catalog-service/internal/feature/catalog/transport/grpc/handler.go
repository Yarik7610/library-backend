package grpc

import (
	"context"

	pb "github.com/Yarik7610/library-backend-common/transport/grpc/microservice/catalog"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/service"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/observability/tracing"
	grpcInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/transport/grpc"
)

type CatalogHandler struct {
	pb.UnimplementedCatalogServiceServer
	config         *config.Config
	logger         *logging.Logger
	catalogService service.CatalogService
}

func NewCatalogHandler(
	config *config.Config,
	logger *logging.Logger,
	catalogService service.CatalogService) *CatalogHandler {
	return &CatalogHandler{
		config:         config,
		logger:         logger,
		catalogService: catalogService,
	}
}

func (h *CatalogHandler) BookCategoryExists(ctx context.Context, req *pb.BookCategoryExistsRequest) (*pb.BookCategoryExistsResponse, error) {
	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.BookCategoryExists")
	defer span.End()

	exists, err := h.catalogService.BookCategoryExists(ctx, req.GetBookCategory())
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Book category exists error", logging.Error(err))
		return nil, grpcInfrastructure.NewError(err)
	}

	return &pb.BookCategoryExistsResponse{Exists: exists}, nil
}
