package http

import (
	"errors"
	"net/http"

	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/dto"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/mapper"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	httpInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/header"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler interface {
	GetBookCategorySubscribedUserEmails(c *gin.Context)
	GetUserSubscribedBookCategories(c *gin.Context)
	SubscribeToBookCategory(c *gin.Context)
	UnsubscribeFromBookCategory(c *gin.Context)
}

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{subscriptionService: subscriptionService}
}

// GetBookCategorySubscribedUserEmails godoc
//
//	@Summary		Get emails of users subscribed to a book category
//	@Description	Returns emails of all users subscribed to the given book category
//	@Tags			internal
//	@Param			categoryName	path	string	true	"Category name"
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/subscriptions/books/categories/{categoryName} [get]
func (h *subscriptionHandler) GetBookCategorySubscribedUserEmails(c *gin.Context) {
	ctx := c.Request.Context()

	bookCategory := c.Param("categoryName")

	emails, err := h.subscriptionService.GetBookCategorySubscribedUserEmails(ctx, bookCategory)
	if err != nil {
		zap.S().Errorf("Get book category subscribed user email error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, emails)
}

// GetUserSubscribedBookCategories godoc
//
//	@Summary		Get book categories the current user is subscribed to
//	@Description	Returns a list of book categories for the user
//	@Tags			subscription
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		string
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		401 {object} 	dto.Error "The token is missing, invalid or expired"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/subscriptions/books/categories [get]
func (h *subscriptionHandler) GetUserSubscribedBookCategories(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := header.GetUserID(c)
	if err != nil {
		httpInfrastructure.RenderError(c, err)
		return
	}

	userSubscribedBookCategories, err := h.subscriptionService.GetUserSubscribedBookCategories(ctx, uint(userID))
	if err != nil {
		zap.S().Errorf("Get user subscribed book categories error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, userSubscribedBookCategories)
}

// SubscribeToBookCategory godoc
//
//	@Summary		Subscribe current user to a book category
//	@Description	Adds the book category to the user's book category subscriptions
//	@Tags			subscription
//	@Param			category	body	dto.SubscribeToBookCategoryRequest	true	"Book category to subscribe"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	dto.UserBookCategory
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		401 {object} 	dto.Error "The token is missing, invalid or expired"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/subscriptions/books/categories [post]
func (h *subscriptionHandler) SubscribeToBookCategory(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := header.GetUserID(c)
	if err != nil {
		httpInfrastructure.RenderError(c, err)
		return
	}

	var subscribeToBookCategoryDTO dto.SubscribeToBookCategoryRequest
	if err := c.ShouldBindJSON(&subscribeToBookCategoryDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	userBookCategoryDomain, err := h.subscriptionService.SubscribeToBookCategory(ctx, uint(userID), subscribeToBookCategoryDTO.BookCategory)
	if err != nil {
		var infrastructureError *errs.Error
		if errors.As(err, &infrastructureError) {
			zap.S().Errorf("Subscribe to book category error: %v", infrastructureError.Cause)
		} else {
			zap.S().Errorf("Subscribe to book category error: %v", err)
		}
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.UserBookCategoryDomainToDTO(userBookCategoryDomain))
}

// UnsubscribeFromBookCategory godoc
//
//	@Summary		Unsubscribe current user from a category
//	@Description	Removes the category from the user's subscriptions
//	@Tags			subscription
//	@Param			categoryName	path	string	true	"Category name to unsubscribe"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	"No content"
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		401 {object} 	dto.Error "The token is missing, invalid or expired"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/subscriptions/books/categories/{categoryName} [delete]
func (h *subscriptionHandler) UnsubscribeFromBookCategory(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := header.GetUserID(c)
	if err != nil {
		httpInfrastructure.RenderError(c, err)
		return
	}

	bookCategory := c.Param("categoryName")

	if err := h.subscriptionService.UnsubscribeFromBookCategory(ctx, uint(userID), bookCategory); err != nil {
		zap.S().Errorf("Unsubscribe from book category error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()
}
