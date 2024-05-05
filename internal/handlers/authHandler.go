package handlers

import (
	// "log"
	"errors"
	"fmt"
	"time"

	// "fmt"
	"net/http"
	"os"
	"waitress-backend/internal/models"

	"waitress-backend/internal/utilities"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func createToken(username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "username": username, 
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 return tokenString, nil
}

func verifyToken(tokenString string) error {
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        return err
    }
    
    if !token.Valid {
        return fmt.Errorf("invalid token")
    }
    
    return nil
}

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

