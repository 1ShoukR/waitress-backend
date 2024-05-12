package handlers

import (
	"fmt"
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
        tx := db.Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", r)})
            }
        }()
        
        var foundUser models.User
        userId := c.PostForm("userId")
        address := c.PostForm("address")
        latitude := c.PostForm("latitude")
        longitude := c.PostForm("longitude")

        if err := tx.Where("user_id = ?", userId).First(&foundUser).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
            return
        }

        la, err := strconv.ParseFloat(latitude, 64)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
            return
        }

        lo, err := strconv.ParseFloat(longitude, 64)
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
            return
        }

        foundUser.Address = &address
        foundUser.Latitude = la
        foundUser.Longitude = lo

        if err := tx.Save(&foundUser).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user location"})
            return
        }

        tx.Commit()
        c.JSON(http.StatusOK, gin.H{"message": "User location updated successfully", "user": foundUser})
    }
}


