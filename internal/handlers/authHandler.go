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
        userAgent := c.PostForm("userAgent")
        var client models.APIClient
        userClient := db.Find(&client, "client_type = ?", userAgent)
        fmt.Println(userClient)
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
        if !utilities.CheckPasswordHash(password, foundUser.PasswordHash, ) {
            c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid login credentials"})
            return
        }
        token, err := createToken((foundUser.Email))
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating token"})
            return
        }
        // If the password is correct, proceed with session handling
        session := sessions.Default(c)
        session.Set("userID", foundUser.UserID)
        session.Set("apiToken", token)
        session.Set("authType", foundUser.AuthType)
        session.Set("loggedIn", true)
        session.Save()

        c.IndentedJSON(http.StatusOK, gin.H{
            "Message": "Login successful", 
            "user": foundUser, 
            "token": token,
            "client": userClient,
        })
    }
}

