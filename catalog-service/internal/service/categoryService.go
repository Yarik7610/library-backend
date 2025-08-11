package service

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type CategoryService interface {
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}
