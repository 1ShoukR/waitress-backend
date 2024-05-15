package handlers

import (
	"net/http"
	"strconv"

	// "waitress-backend/internal/handlers"
	"waitress-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		if email == "" || password == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid login credentials"})
			return
		}
		// salt := handlers.GenerateSalt(16)
		// hashedPassword := handlers.HashPassword()
	}
	// logic to create a new user
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		results := db.Find(&users)
		if results.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": results.Error.Error()})
		}
		c.IndentedJSON(http.StatusOK, users)
	}
}

func UpdateUserLocation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.User
		userId := c.PostForm("userId")
		address := c.PostForm("address")
		latitudeStr := c.PostForm("latitude")
		longitudeStr := c.PostForm("longitude")

		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}

		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		if err := db.Where("user_id = ?", userId).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Now use the method on the foundUser object
		if err := foundUser.UpdateLocation(db, latitude, longitude, address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User location updated successfully", "user": foundUser})
	}
}
