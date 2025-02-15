package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/internal/services"
)

type UserInfoHandler struct {
	service *services.UserInfoService
}
func NewUserInfoHandler(service *services.UserInfoService) *UserInfoHandler {
	return &UserInfoHandler{service}
}

func (h *UserInfoHandler) GetUserInfo(c *gin.Context) {
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

	userDetail, err := h.service.GetUserInfo(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userDetail)

}
