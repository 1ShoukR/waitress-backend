package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"waitress-backend/internal/models"
)


func CreateUser(c *gin.Context) {
    response := "Test User Created"
    c.IndentedJSON(http.StatusOK, response)
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

