package handlers

import (
	// "log"
	"errors"
	// "fmt"
	"net/http"

	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"waitress-backend/internal/utilities"
)

func Login(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
    return func(c *gin.Context) {
        email := c.PostForm("email")
        password := c.PostForm("password")
        if email == "" {
            c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The email and/or password is incorrect"})
            return
        }

        var foundUser models.User
        result := db.Where("email = ?", email).First(&foundUser)
        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
            }
            return
        }

        // Assume foundUser.Salt and foundUser.PasswordHash store the salt and hashed password
        if !utilities.VerifyPassword(foundUser.PasswordHash, password, []byte(foundUser.Salt)) {
            c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid login credentials"})
            return
        }

        // If the password is correct, proceed with session handling
        session := sessions.Default(c)
        var count int
        v := session.Get("Count")
        if v == nil {
            count = 0
        } else {
            count = v.(int) + 1
        }
        session.Set("Count", count)
        session.Save()

        c.IndentedJSON(http.StatusOK, gin.H{"Message": "Login successful", "Count": count})
    }
}

