package postgres

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func NewError(err error) *errs.Error {
	const postgresUniqueViolationCode = "23505"

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewEntityNotFoundError()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == postgresUniqueViolationCode {
			return errs.NewEntityAlreadyExistsError().WithCause(err)
		}
	}
	return errs.NewInternalServerError().WithCause(err)
}
