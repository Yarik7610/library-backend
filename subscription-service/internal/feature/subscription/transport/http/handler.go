package http

import (
	"net/http"

	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/dto"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/mapper"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/logging"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/observability/tracing"
	httpInfrastructure "github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/transport/http/header"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler interface {
	GetBookCategorySubscribedUserEmails(c *gin.Context)
	GetUserSubscribedBookCategories(c *gin.Context)
	SubscribeToBookCategory(c *gin.Context)
	UnsubscribeFromBookCategory(c *gin.Context)
}

type subscriptionHandler struct {
	config              *config.Config
	logger              *logging.Logger
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(
	config *config.Config,
	logger *logging.Logger,
	subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{
		config:              config,
		logger:              logger,
		subscriptionService: subscriptionService,
	}
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

	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.GetBookCategorySubscribedUserEmails")
	defer span.End()

	emails, err := h.subscriptionService.GetBookCategorySubscribedUserEmails(ctx, bookCategory)
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Get book category subscribed user email error", logging.Error(err))
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

	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.GetUserSubscribedBookCategories")
	defer span.End()

	userSubscribedBookCategories, err := h.subscriptionService.GetUserSubscribedBookCategories(ctx, uint(userID))
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Get user subscribed book categories error", logging.Error(err))
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

	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.SubscribeToBookCategory")
	defer span.End()

	userBookCategoryDomain, err := h.subscriptionService.SubscribeToBookCategory(ctx, uint(userID), subscribeToBookCategoryDTO.BookCategory)
	if err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Subscribe to book category error", logging.Error(err))
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

	ctx, span := tracing.Span(ctx, h.config.ServiceName, "service.UnsubscribeFromBookCategory")
	defer span.End()

	if err := h.subscriptionService.UnsubscribeFromBookCategory(ctx, uint(userID), bookCategory); err != nil {
		tracing.Error(span, err)
		h.logger.Error(ctx, "Unsubscribe from book category error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()
}
