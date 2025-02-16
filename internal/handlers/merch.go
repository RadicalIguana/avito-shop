package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/RadicalIguana/avito-shop/internal/services"
)

type MerchHandler struct {
	service *services.MerchService
}

func NewMerchHandler(service *services.MerchService) *MerchHandler {
    return &MerchHandler{service}
}

func (h *MerchHandler) PurchaseItem(c *gin.Context) {
	itemName := c.Param("item")
	ctx := c.Request.Context()

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

	if err := h.service.PurchaseItem(ctx, userID, itemName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item purchased successfully",
	})
}