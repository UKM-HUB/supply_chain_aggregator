package main

import (
	"golang-api/config"
	"golang-api/handlers"
	"golang-api/models"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()

	config.DB.AutoMigrate(
		&models.User{},
		&models.UMKM{},
		&models.Factory{},
		&models.Transaction{},
	)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/users", handlers.GetUsers)
		api.GET("/umkm", handlers.GetUMKM)
		api.GET("/factories", handlers.GetFactories)
		api.GET("/transactions", handlers.GetTransactions)
	}

	r.Run(":8080")
}