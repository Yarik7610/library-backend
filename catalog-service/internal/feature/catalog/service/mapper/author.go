package mapper

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
)

func AuthorModelToDomain(authorModel *model.Author) domain.Author {
	return domain.Author{
		ID:       authorModel.ID,
		Fullname: authorModel.Fullname,
	}
}
