package controller

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
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

}

func (c *subscriptionController) SubscribeCategory(ctx *gin.Context) {

}

func (c *subscriptionController) UnsubscribeCategory(ctx *gin.Context) {

}
