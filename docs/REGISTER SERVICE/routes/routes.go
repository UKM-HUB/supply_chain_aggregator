package routes

import (
	"register-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")
	{
		api.POST("/register", handlers.Register)
		api.POST("/register/verify", handlers.VerifyOTP)
	}
}