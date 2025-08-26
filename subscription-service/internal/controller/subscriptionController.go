package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
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
