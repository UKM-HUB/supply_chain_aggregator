package main

import (
	"os"

	"xendit-va/config"
	"xendit-va/database"
	"xendit-va/routes"
	"xendit-va/services"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	database.ConnectDB()

	services.InitRabbitMQ()

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")

	r.Run(":" + port)
}