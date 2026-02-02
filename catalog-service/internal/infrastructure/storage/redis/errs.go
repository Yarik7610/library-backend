package redis

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	"github.com/redis/go-redis/v9"
)

func NewError(err error, entityName string) *errs.Error {
	if errors.Is(err, redis.Nil) {
		return errs.NewEntityNotFoundError(entityName)
	}
	return errs.NewInternalServerError().WithCause(err)
}
