package handlers

import (
	"fmt"
	"net/http"
	"waitress-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAdminData(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		fmt.Printf("user id: %s\n", userId) // Add newline for better logging
		
		returnData := make(map[string]interface{})
		var user models.User
		
		if err := db.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		
		returnData["user"] = user
		c.JSON(http.StatusOK, gin.H{"success": returnData})
	}
}