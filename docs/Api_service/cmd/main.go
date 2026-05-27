package main

import (
	"os"
	"project/config"
	"project/models"
	"project/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	config.ConnectDB()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":" + os.Getenv("APP_PORT"))
}