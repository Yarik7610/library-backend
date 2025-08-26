package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
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
	exists, err := s.categoryExists(category)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return nil, custom.NewErr(http.StatusBadRequest, "can't subscribe on unknown category")
	}

	subscribedCategory := model.UserCategory{
		UserID:   userID,
		Category: category,
	}
	err = s.userCategoryRepository.SubscribeCategory(&subscribedCategory)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return &subscribedCategory, nil
}

func (s *catalogService) UnsubscribeCategory(userID uint, category string) *custom.Err {
	exists, err := s.categoryExists(category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return custom.NewErr(http.StatusBadRequest, "can't unsubscribe from unknown category")
	}

	subscribedCategory, err := s.userCategoryRepository.FindSubscribedCategory(userID, category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if subscribedCategory == nil {
		return custom.NewErr(http.StatusBadRequest, "didn't find such category in your subscriptions")
	}

	err = s.userCategoryRepository.UnsubscribeCategory(userID, category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) categoryExists(category string) (bool, error) {
	resp, err := http.Get(sharedconstants.CATALOG_MICROSERVICE_SOCKET + sharedconstants.CATALOG_ROUTE + sharedconstants.CATEGORIES_ROUTE)
	if err != nil {
		return false, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, custom.NewErr(http.StatusInternalServerError, fmt.Sprintf("catalog microservice return status code: %d", resp.StatusCode))
	}

	var categories []string
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return false, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	category = strings.ToLower(category)
	return slices.Contains(categories, category), nil
}
