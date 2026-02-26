package grpc

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(err error) error {
	var infrastructureError *errs.Error
	if errors.As(err, &infrastructureError) {
		return status.Errorf(getGRPCCode(infrastructureError.Code), "%s", infrastructureError.Message)
	}
	return status.Errorf(codes.Internal, "Internal server error")
}

func getGRPCCode(errorCode errs.Code) codes.Code {
	errorCodesToGRPCCodes := map[errs.Code]codes.Code{
		errs.CodeNotFound:      codes.NotFound,
		errs.CodeAlreadyExists: codes.AlreadyExists,
		errs.CodeBadRequest:    codes.InvalidArgument,
		errs.CodeInternal:      codes.Internal,
	}
	if code, exists := errorCodesToGRPCCodes[errorCode]; exists {
		return code
	}
	return codes.Internal
}
