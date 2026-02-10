package http

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/transport/http/header"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/service"
	"github.com/Yarik7610/library-backend/subscription-service/internal/feature/subscription/transport/http/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler interface {
	GetCategorySubscribersEmails(c *gin.Context)
	GetUserBookCategories(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{subscriptionService: subscriptionService}
}

// GetCategorySubscribersEmails godoc
//
//	@Summary		Get emails of users subscribed to a category
//	@Description	Returns emails of all users subscribed to the given category
//	@Tags			internal
//	@Param			categoryName	path	string	true	"Category name"
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/books/categories/{categoryName} [get]
func (h *subscriptionHandler) GetCategorySubscribersEmails(c *gin.Context) {
	category := c.Param("categoryName")

	emails, err := h.subscriptionService.GetCategorySubscribersEmails(category)
	if err != nil {
		zap.S().Errorf("Get category subscribers IDs error: %v\n", err)
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, emails)
}

// GetUserBookCategories godoc
//
//	@Summary		Get categories the current user is subscribed to
//	@Description	Returns a list of categories for the user
//	@Tags			subscription
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/books/categories [get]
func (h *subscriptionHandler) GetUserBookCategories(c *gin.Context) {
	userIDString := c.GetHeader(header.USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribedCategories, customErr := h.subscriptionService.GetUserBookCategories(uint(userID))
	if customErr != nil {
		zap.S().Errorf("Get subscribed categories error: %v\n", customErr)
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	c.JSON(http.StatusOK, subscribedCategories)
}

// Create godoc
//
//	@Summary		Subscribe current user to a category
//	@Description	Adds the category to the user's subscriptions
//	@Tags			subscription
//	@Param			category	body	dto.Create	true	"Category to subscribe"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/books/categories [post]
func (h *subscriptionHandler) Create(c *gin.Context) {
	userIDString := c.GetHeader(header.USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscribeCategoryDTO dto.Create
	if err := c.ShouldBindJSON(&subscribeCategoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribedCategory, customErr := h.subscriptionService.Create(uint(userID), subscribeCategoryDTO.Category)
	if customErr != nil {
		zap.S().Errorf("Subscribe category error: %v\n", customErr)
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	c.JSON(http.StatusOK, subscribedCategory)
}

// Delete godoc
//
//	@Summary		Unsubscribe current user from a category
//	@Description	Removes the category from the user's subscriptions
//	@Tags			subscription
//	@Param			categoryName	path	string	true	"Category name to unsubscribe"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/books/categories/{categoryName} [delete]
func (h *subscriptionHandler) Delete(c *gin.Context) {
	userIDString := c.GetHeader(header.USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := c.Param("categoryName")

	customErr := h.subscriptionService.Delete(uint(userID), category)
	if customErr != nil {
		zap.S().Errorf("Unsubscribe category error: %v\n", customErr)
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	c.Status(http.StatusNoContent)
	c.Abort()
}
