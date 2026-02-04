package redis

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	"github.com/redis/go-redis/v9"
)

func NewError(err error) *errs.Error {
	return errs.NewInternalServerError().WithCause(err)
}

func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
