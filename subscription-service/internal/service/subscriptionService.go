package service

import (
	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type SubscriptionService interface {
	GetSubscribedCategories() ([]string, *custom.Err)
}

type catalogService struct {
	userCategoryRepository repository.UserCategoryRepository
}

func NewSubscriptionService(userCategoryRepository repository.UserCategoryRepository) SubscriptionService {
	return &catalogService{userCategoryRepository: userCategoryRepository}
}

func (s *catalogService) GetSubscribedCategories() ([]string, *custom.Err) {
	return nil, nil
}
