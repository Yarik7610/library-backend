package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/subscription-service/internal/model"
	"github.com/Yarik7610/library-backend/subscription-service/internal/repository"
)

type SubscriptionService interface {
	GetCategorySubscribersEmails(category string) ([]string, *custom.Err)
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

func (s *catalogService) GetCategorySubscribersEmails(category string) ([]string, *custom.Err) {
	subscribersIDs, err := s.userCategoryRepository.GetCategorySubscribersIDs(category)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	emails, customErr := s.getEmailsByUserIDs(subscribersIDs)
	if customErr != nil {
		return nil, customErr
	}
	return emails, nil
}

func (s *catalogService) GetSubscribedCategories(userID uint) ([]string, *custom.Err) {
	subscribedCategories, err := s.userCategoryRepository.GetSubscribedCategories(userID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return subscribedCategories, nil
}

func (s *catalogService) SubscribeCategory(userID uint, category string) (*model.UserCategory, *custom.Err) {
	exists, customErr := s.categoryExists(category)
	if customErr != nil {
		return nil, customErr
	}
	if !exists {
		return nil, custom.NewErr(http.StatusBadRequest, "can't subscribe on unknown category")
	}

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
	exists, customErr := s.categoryExists(category)
	if customErr != nil {
		return customErr
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

func (s *catalogService) categoryExists(category string) (bool, *custom.Err) {
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

func (s *catalogService) getEmailsByUserIDs(userIDs []uint) ([]string, *custom.Err) {
	req, err := http.NewRequest("GET", sharedconstants.USER_MICROSERVICE_SOCKET+sharedconstants.EMAILS_ROUTE, nil)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	q := req.URL.Query()
	for _, userID := range userIDs {
		q.Add("ids", strconv.FormatUint(uint64(userID), 10))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, custom.NewErr(http.StatusInternalServerError, fmt.Sprintf("user microservice return status code: %d", resp.StatusCode))
	}

	var emails []string
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return emails, nil
}
