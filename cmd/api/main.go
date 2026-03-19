// Package main
package main

import (
	"log"
	"os"

	"github.com/fbveronez/go-test-case-api/internal/db"
	"github.com/fbveronez/go-test-case-api/internal/handlers"
	"github.com/fbveronez/go-test-case-api/internal/repository"
	"github.com/fbveronez/go-test-case-api/internal/service"

	_ "github.com/fbveronez/go-test-case-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Test Case API
// @version 1.0
// @description Test API built in Go + Gin
// @host localhost:8080
// @BasePath /

func main() {
	db.Connect()
	log.Println("Database connected!")

	router := gin.Default()

	accountRepo := repository.NewAccountRepository(db.DB)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountService)

	transactionRepo := repository.NewTransactionRepository(db.DB)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService, accountService)

	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:id", accountHandler.GetAccountByID)
	router.DELETE("/accounts/:id", accountHandler.DeleteAccount)

	router.POST("/transactions", transactionHandler.CreateTransaction)
	router.GET("/transactions/:id", transactionHandler.GetTransactionsByAccountID)
	router.PUT("/transactions/:id", transactionHandler.UpdateTransaction)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/ui", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
