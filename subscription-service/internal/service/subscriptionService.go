package service

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	GetSubscribedCategories(userID uint) ([]string, *custom.Err)
	SubscribeCategory(userID uint, category string) (*model.UserCategory, *custom.Err)
	UnsubscribeCategory(userID uint, category string) *custom.Err
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

func (s *catalogService) SubscribeCategory(userID uint, category string) (*model.UserCategory, *custom.Err) {
	subscribedCategory := model.UserCategory{
		UserID:   userID,
		Category: category,
	}
	err := s.userCategoryRepository.SubscribeCategory(&subscribedCategory)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return &subscribedCategory, nil
}

func (s *catalogService) UnsubscribeCategory(userID uint, category string) *custom.Err {
	subscribedCategory, err := s.userCategoryRepository.FindSubscribedCategory(userID, category)
	zap.S().Debug(subscribedCategory)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if subscribedCategory == nil {
		return custom.NewErr(http.StatusBadRequest, "didn't find such category in subscribed")
	}

	err = s.userCategoryRepository.UnsubscribeCategory(userID, category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}
