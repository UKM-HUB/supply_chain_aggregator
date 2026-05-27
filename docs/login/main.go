package main

import "github.com/gin-gonic/gin"

func main() {

	ConnectDB()

	r := gin.Default()

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
		auth.POST("/refresh", Refresh)
		auth.POST("/logout", Logout)
	}

	protected := api.Group("/user")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Protected route",
			})
		})
	}

	r.Run(":8080")
}