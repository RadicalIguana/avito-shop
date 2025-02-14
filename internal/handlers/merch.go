package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	// "github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/services"
)

// TODO: PurchaseItemHandler?
type MerchHandler struct {
	service *services.MerchService
}

func NewMerchHandler(service *services.MerchService) *MerchHandler {
    return &MerchHandler{service}
}

// TODO: Определить единое название
// TODO: Что такое gin.HandlerFunc?
func (h *MerchHandler) PurchaseItem(c *gin.Context) {
	// var req models.BuyRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// TODO: Извлечение ID покупателя из JWT
	userID := "1"

	itemName := c.Param("item")
	ctx := c.Request.Context()

	if err := h.service.PurchaseItem(ctx, userID, itemName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item purchased successfully",
	})
}