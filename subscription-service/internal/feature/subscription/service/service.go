package service

import (
	"context"

	"github.com/Yarik7610/library-backend/subscription-service/internal/domain"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/repository/postgres/model"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service/mapper"

	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/catalog"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/microservice/user"
)

type SubscriptionService interface {
	GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error)
	GetUserSubscribedBookCategories(ctx context.Context, userID uint) ([]string, error)
	SubscribeToBookCategory(ctx context.Context, userID uint, bookCategory string) (*domain.UserBookCategory, error)
	UnsubscribeFromBookCategory(ctx context.Context, userID uint, bookCategory string) error
}

type subscriptionService struct {
	userBookCategorySubscriptionRepository postgres.UserBookCategorySubscriptionRepository
	catalogMicroserviceClient              catalog.Client
	userMicroserviceClient                 user.Client
}

func NewSubscriptionService(
	userBookCategorySubscriptionRepository postgres.UserBookCategorySubscriptionRepository,
	catalogMicroserviceClient catalog.Client,
	userMicroserviceClient user.Client,
) SubscriptionService {
	return &subscriptionService{
		userBookCategorySubscriptionRepository: userBookCategorySubscriptionRepository,
		catalogMicroserviceClient:              catalogMicroserviceClient,
		userMicroserviceClient:                 userMicroserviceClient,
	}
}

func (s *subscriptionService) GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error) {
	userIDs, err := s.userBookCategorySubscriptionRepository.GetSubscriptionUserIDs(ctx, bookCategory)
	if err != nil {
		return nil, err
	}

	emails, err := s.userMicroserviceClient.GetEmailsByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (s *subscriptionService) GetUserSubscribedBookCategories(ctx context.Context, userID uint) ([]string, error) {
	subscribedUserBookCategories, err := s.userBookCategorySubscriptionRepository.GetUserSubscribedBookCategories(ctx, userID)
	if err != nil {
		return nil, err
	}
	return subscribedUserBookCategories, nil
}

func (s *subscriptionService) SubscribeToBookCategory(ctx context.Context, userID uint, bookCategory string) (*domain.UserBookCategory, error) {
	exists, err := s.catalogMicroserviceClient.BookCategoryExists(ctx, bookCategory)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.NewEntityNotFoundError("Book category")
	}

	userBookCategoryModel := model.UserBookCategory{
		UserID:       userID,
		BookCategory: bookCategory,
	}
	if err := s.userBookCategorySubscriptionRepository.Create(ctx, &userBookCategoryModel); err != nil {
		return nil, err
	}

	userBookCategoryDomain := mapper.UserBookCategoryToDomain(&userBookCategoryModel)
	return &userBookCategoryDomain, nil
}

func (s *subscriptionService) UnsubscribeFromBookCategory(ctx context.Context, userID uint, bookCategory string) error {
	exists, err := s.catalogMicroserviceClient.BookCategoryExists(ctx, bookCategory)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewEntityNotFoundError("Book category")
	}

	return s.userBookCategorySubscriptionRepository.Delete(ctx, userID, bookCategory)
}
