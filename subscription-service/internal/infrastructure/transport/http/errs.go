package http

import (
	"errors"
	"net/http"

	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/dto"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	"github.com/gin-gonic/gin"
)

func RenderError(c *gin.Context, err error) {
	var infrastructureError *errs.Error
	if errors.As(err, &infrastructureError) {
		c.JSON(getHTTPStatus(infrastructureError.Code), dto.Error{Error: infrastructureError.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, dto.Error{Error: "Internal server error"})
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
