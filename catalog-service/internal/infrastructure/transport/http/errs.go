package http

import (
	"errors"
	"net/http"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/transport/http/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	"github.com/gin-gonic/gin"
)

func RenderError(ctx *gin.Context, err error) {
	var infrastructureError *errs.Error
	if errors.As(err, &infrastructureError) {
		ctx.JSON(getHTTPStatus(infrastructureError.Code), dto.Error{Error: infrastructureError.Message})
		return
	}
	ctx.JSON(http.StatusInternalServerError, dto.Error{Error: "Internal server error"})
}

func getHTTPStatus(errorCode errs.Code) int {
	errorCodesToHTTPStatuses := map[errs.Code]int{
		errs.CodeNotFound:      http.StatusNotFound,
		errs.CodeAlreadyExists: http.StatusConflict,
		errs.CodeBadRequest:    http.StatusBadRequest,
		errs.CodeInternal:      http.StatusInternalServerError,
	}

	if status, exists := errorCodesToHTTPStatuses[errorCode]; exists {
		return status
	}
	return http.StatusInternalServerError
}
