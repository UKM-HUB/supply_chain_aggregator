package routes

import (
	"project/controllers"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	{
		api.POST("/auth/register", controllers.Register)
		api.POST("/auth/login", controllers.Login)

		api.GET("/users", middleware.AuthMiddleware(), controllers.GetUsers)

		api.POST("/transactions", middleware.AuthMiddleware(), controllers.CreateTransaction)
		api.GET("/transactions", middleware.AuthMiddleware(), controllers.GetTransactions)

		api.GET("/reports", middleware.AuthMiddleware(), controllers.GetReports)

		api.GET("/nearby", controllers.Nearby)

		api.GET("/communication", controllers.Communication)

		api.GET("/gateway", controllers.Gateway)
	}
}