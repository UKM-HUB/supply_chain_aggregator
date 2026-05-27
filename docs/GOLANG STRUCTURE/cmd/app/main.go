package main

import (
    "log"
    "os"
    "project/internal/database"
    delivery "project/internal/delivery/http"
    "project/internal/delivery/http/handler"
    "project/internal/entity"
    "project/internal/helper"
    "project/internal/repository"
    "project/internal/usecase"

    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
)

func main() {

    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    db, err := database.ConnectDB()
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(
        &entity.User{},
        &entity.Transaction{},
    )

    helper.InitRabbitMQ(os.Getenv("RABBITMQ_URL"))

    userRepo := repository.NewUserRepository(db)
    txRepo := repository.NewTransactionRepository(db)

    authUsecase := usecase.NewAuthUsecase(userRepo)
    txUsecase := usecase.NewTransactionUsecase(txRepo)

    authHandler := handler.NewAuthHandler(authUsecase)
    txHandler := handler.NewTransactionHandler(txUsecase)

    e := echo.New()

    delivery.RegisterRoutes(
        e,
        authHandler,
        txHandler,
    )

    e.Logger.Fatal(e.Start(":8080"))
}