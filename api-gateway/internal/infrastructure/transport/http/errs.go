package http

import (
	"errors"
	"net/http"

	"github.com/Yarik7610/library-backend/api-gateway/internal/app/dto"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/errs"
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
		errs.CodeUnautorized: http.StatusUnauthorized,
		errs.CodeForbidden:   http.StatusForbidden,
		errs.CodeInternal:    http.StatusInternalServerError,
	}

	if status, exists := errorCodesToHTTPStatuses[errorCode]; exists {
		return status
	}
	return http.StatusInternalServerError
}
