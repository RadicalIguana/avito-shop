package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/RadicalIguana/avito-shop/internal/database"
	"github.com/RadicalIguana/avito-shop/internal/handlers"
	"github.com/RadicalIguana/avito-shop/internal/middlewares"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
	"github.com/RadicalIguana/avito-shop/internal/services"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// TODO: Для чего такие определения?
	coinRepo := repositories.NewCoinRepository(database.DB)
	coinService := services.NewCoinService(coinRepo)
	coinHandler := handlers.NewCoinHandler(coinService)

	merchRepo := repositories.NewMerchRepository(database.DB)
	merchService := services.NewMerchService(merchRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	userRepo := repositories.NewUserInfoRepository(database.DB)
	userService := services.NewUserInfoService(userRepo)
	userHandler := handlers.NewUserInfoHandler(userService)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.RedirectTrailingSlash = false
	r.RemoveExtraSlash = true

	r.POST("/api/auth", handlers.AuthHandler)

	r.Use(middlewares.AuthMiddleware())

	r.POST("/api/sendCoin", coinHandler.SendCoins)
	r.GET("/api/buy/:item", merchHandler.PurchaseItem)
	r.GET("/api/info", userHandler.GetUserInfo)

	log.Fatal(r.Run("0.0.0.0:8080"))
}
