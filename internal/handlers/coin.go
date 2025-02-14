package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/services"
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

	// Извлечение ID отправителя из JWT
	// claims := r.Context().Value("userClaims").(map[string]interface{})
    // fromUserID := claims["userID"].(string)

    // TODO: Изменить получение ID отправителя 
	fromUserID := "1"
	
	if err := h.service.TransferCoins(c.Request.Context(), fromUserID, req.ToUser, req.Amount); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "message": "Transfer successful",
    })
}