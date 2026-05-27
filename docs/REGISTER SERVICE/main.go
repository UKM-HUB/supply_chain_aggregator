package main

import (
	"register-api/config"
	"register-api/models"
	"register-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()

	config.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":8080")
}