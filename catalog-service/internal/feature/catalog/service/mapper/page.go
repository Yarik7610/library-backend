package mapper

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
)

func PageModelToDomain(pageModel *model.Page) domain.Page {
	return domain.Page{
		ID:      pageModel.ID,
		Number:  pageModel.Number,
		Content: pageModel.Content,
	}
}
