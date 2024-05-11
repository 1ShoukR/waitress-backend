package handlers

import (
	"fmt"
	"net/http"

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
        // logic to update user location
        var foundUser models.User
        userId := c.PostForm("userId")
        if err := db.Where("user_id = ?", userId).First(&foundUser).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
            return
        }
        fmt.Println(foundUser)
    }
}

