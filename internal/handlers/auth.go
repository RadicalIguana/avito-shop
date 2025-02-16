package handlers

import (
	"context"
	"net/http"

	"github.com/RadicalIguana/avito-shop/internal/database"
	"github.com/RadicalIguana/avito-shop/internal/models"
	"github.com/RadicalIguana/avito-shop/utils"
	"github.com/gin-gonic/gin"
)

func AuthHandler(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	var user models.User
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT id, password FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.Id, &user.Password)

	if err != nil {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

		err = database.DB.QueryRow(
			context.Background(),
            "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
            req.Username, hashedPassword,
		).Scan(&user.Id)

		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
	} else {
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
            return
		}
	}

	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"token": token})
}