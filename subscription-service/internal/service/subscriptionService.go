package service

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
)

type SubscriptionService interface {
	GetSubscribedCategories(userID uint) ([]string, *custom.Err)
}

type catalogService struct {
	userCategoryRepository repository.UserCategoryRepository
}

func NewSubscriptionService(userCategoryRepository repository.UserCategoryRepository) SubscriptionService {
	return &catalogService{userCategoryRepository: userCategoryRepository}
}

func (s *catalogService) GetSubscribedCategories(userID uint) ([]string, *custom.Err) {
	subscribedCategories, err := s.userCategoryRepository.GetSubscribedCategories(userID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return subscribedCategories, nil
}
