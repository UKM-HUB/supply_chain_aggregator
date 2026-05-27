package routes

import (
	"xendit-va/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	{
		api.POST("/gateway/create-va", handlers.CreateVA)
		api.POST("/webhooks/xendit", handlers.XenditWebhookHandler)
	}
}