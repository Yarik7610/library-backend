package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Yarik7610/library-backend-common/custom"
	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres/model"
)

type SubscriptionService interface {
	GetCategorySubscribersEmails(category string) ([]string, *custom.Err)
	GetUserBookCategories(userID uint) ([]string, *custom.Err)
	Create(userID uint, category string) (*model.UserBookCategory, *custom.Err)
	Delete(userID uint, category string) *custom.Err
}

type catalogService struct {
	userBookCategoryRepository postgres.UserBookCategoryRepository
}

func NewSubscriptionService(userBookCategoryRepository postgres.UserBookCategoryRepository) SubscriptionService {
	return &catalogService{userBookCategoryRepository: userBookCategoryRepository}
}

func (s *catalogService) GetCategorySubscribersEmails(category string) ([]string, *custom.Err) {
	subscribersIDs, err := s.userBookCategoryRepository.GetBookCategoryUserIDs(category)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}

	emails, customErr := s.getEmailsByUserIDs(subscribersIDs)
	if customErr != nil {
		return nil, customErr
	}
	return emails, nil
}

func (s *catalogService) GetUserBookCategories(userID uint) ([]string, *custom.Err) {
	subscribedCategories, err := s.userBookCategoryRepository.GetUserBookCategories(userID)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return subscribedCategories, nil
}

func (s *catalogService) Create(userID uint, category string) (*model.UserBookCategory, *custom.Err) {
	exists, customErr := s.categoryExists(category)
	if customErr != nil {
		return nil, customErr
	}
	if !exists {
		return nil, custom.NewErr(http.StatusBadRequest, "can't subscribe on unknown category")
	}

	subscribedCategory := model.UserBookCategory{
		UserID:   userID,
		Category: category,
	}
	err := s.userBookCategoryRepository.Create(&subscribedCategory)
	if err != nil {
		return nil, custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return &subscribedCategory, nil
}

func (s *catalogService) Delete(userID uint, category string) *custom.Err {
	exists, customErr := s.categoryExists(category)
	if customErr != nil {
		return customErr
	}
	if !exists {
		return custom.NewErr(http.StatusBadRequest, "can't unsubscribe from unknown category")
	}

	subscribedCategory, err := s.userBookCategoryRepository.FindUserBookCategory(userID, category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	if subscribedCategory == nil {
		return custom.NewErr(http.StatusBadRequest, "didn't find such category in your subscriptions")
	}

	err = s.userBookCategoryRepository.Delete(userID, category)
	if err != nil {
		return custom.NewErr(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *catalogService) categoryExists(category string) (bool, *custom.Err) {
	resp, err := http.Get(microservice.CATALOG_ADDRESS + route.CATALOG + route.CATEGORIES)
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
	req, err := http.NewRequest("GET", microservice.USER_ADDRESS+route.EMAILS, nil)
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
