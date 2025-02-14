package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/RadicalIguana/avito-shop/internal/database"
	"github.com/RadicalIguana/avito-shop/internal/handlers"
	"github.com/RadicalIguana/avito-shop/internal/services"
	"github.com/RadicalIguana/avito-shop/internal/repositories"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: %w", err)
	}
	defer db.Close()

	// TODO: Для чего такие определения?
	coinRepo := repositories.NewCoinRepository(db)
	coinService := services.NewCoinService(coinRepo)
	coinHandler := handlers.NewCoinHandler(coinService)

    merchRepo := repositories.NewMerchRepository(db)
	merchService := services.NewMerchService(merchRepo)
	merchHandler := handlers.NewMerchHandler(merchService)

	r := gin.Default()
	r.POST("/api/sendCoin", coinHandler.SendCoins)
	r.GET("/api/buy/:item", merchHandler.PurchaseItem)

	log.Fatal(r.Run(":8080")) // listen and serve on 0.0.0.0:8080
}
