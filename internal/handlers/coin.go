package handlers

import (
	"net/http"

	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/services"
	"github.com/gin-gonic/gin"
)

type CoinHandler struct {
	service *services.CoinService
}

func NewCoinHandler(service *services.CoinService) *CoinHandler {
    return &CoinHandler{service}
}

func (h *CoinHandler) SendCoins(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	fromUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	floatID, ok := fromUserID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	userID := int(floatID)

	if err := h.service.TransferCoins(c.Request.Context(), userID, req.ToUser, req.Amount); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "message": "Transfer successful",
    })
}