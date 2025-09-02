package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/subscription-service/internal/dto"
	"github.com/Yarik7610/library-backend/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionController interface {
	GetCategorySubscribersEmails(ctx *gin.Context)
	GetSubscribedCategories(ctx *gin.Context)
	SubscribeCategory(ctx *gin.Context)
	UnsubscribeCategory(ctx *gin.Context)
}

type subscriptionController struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionController(subscriptionService service.SubscriptionService) SubscriptionController {
	return &subscriptionController{subscriptionService: subscriptionService}
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
//	@Router			/subscriptions/categories/{categoryName} [get]
func (c *subscriptionController) GetCategorySubscribersEmails(ctx *gin.Context) {
	category := ctx.Param("categoryName")

	emails, err := c.subscriptionService.GetCategorySubscribersEmails(category)
	if err != nil {
		zap.S().Errorf("Get category subscribers IDs error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, emails)
}

// GetSubscribedCategories godoc
//
//	@Summary		Get categories the current user is subscribed to
//	@Description	Returns a list of categories for the user
//	@Tags			subscription
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/categories [get]
func (c *subscriptionController) GetSubscribedCategories(ctx *gin.Context) {
	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribedCategories, customErr := c.subscriptionService.GetSubscribedCategories(uint(userID))
	if customErr != nil {
		zap.S().Errorf("Get subscribed categories error: %v\n", customErr)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, subscribedCategories)
}

// SubscribeCategory godoc
//
//	@Summary		Subscribe current user to a category
//	@Description	Adds the category to the user's subscriptions
//	@Tags			subscription
//	@Param			category	body	dto.SubscribeCategory	true	"Category to subscribe"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/categories [post]
func (c *subscriptionController) SubscribeCategory(ctx *gin.Context) {
	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscribeCategoryDTO dto.SubscribeCategory
	if err := ctx.ShouldBindJSON(&subscribeCategoryDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribedCategory, customErr := c.subscriptionService.SubscribeCategory(uint(userID), subscribeCategoryDTO.Category)
	if customErr != nil {
		zap.S().Errorf("Subscribe category error: %v\n", customErr)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, subscribedCategory)
}

// UnsubscribeCategory godoc
//
//	@Summary		Unsubscribe current user from a category
//	@Description	Removes the category from the user's subscriptions
//	@Tags			subscription
//	@Param			categoryName	path	string	true	"Category name to unsubscribe"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/subscriptions/categories/{categoryName} [delete]
func (c *subscriptionController) UnsubscribeCategory(ctx *gin.Context) {
	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := ctx.Param("categoryName")

	customErr := c.subscriptionService.UnsubscribeCategory(uint(userID), category)
	if customErr != nil {
		zap.S().Errorf("Unsubscribe category error: %v\n", customErr)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.Status(http.StatusNoContent)
	ctx.Abort()
}
