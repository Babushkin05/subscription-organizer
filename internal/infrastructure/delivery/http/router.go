package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *SubscriptionHandler) {
	s := r.Group("/subscriptions")
	{
		s.POST("", handler.CreateSubscription)
		s.GET("", handler.ListSubscriptions)
		s.GET("/cost", handler.CalculateTotalCost)
		s.GET("/:id", handler.GetSubscription)
		s.PUT("/:id", handler.UpdateSubscription)
		s.DELETE("/:id", handler.DeleteSubscription)
	}
}
