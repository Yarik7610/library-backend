package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionController interface {
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

}

func (c *subscriptionController) UnsubscribeCategory(ctx *gin.Context) {

}
