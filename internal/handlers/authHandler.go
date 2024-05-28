// This file contains the handlers for the authentication endpoints
//
// The handlers here are as follows:
// - Login
// - Logout
// - createToken
// - verifyToken
package handlers

import (
	// "log"
	"errors"
	"fmt"
	"log"
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

// Here is how you can access the JWT_SECRET environment variable
var secretKey = []byte(os.Getenv("JWT_SECRET"))

// Function to create an authentication token for a user
func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Function to verify the authenticity of a token
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

// Logout function to clear the session
func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}

// Login function to authenticate a user
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
		result := db.Preload("Entity").Where("email = ?", email).First(&foundUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			}
			return
		}

		// Assume foundUser.Salt and foundUser.PasswordHash store the salt and hashed password
		if !utilities.CheckPasswordHash(password, foundUser.PasswordHash) {
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
		session.Set("clientType", userAgent)
		session.Set("user", foundUser)
		session.Set("loggedIn", true)
		if err := session.Save(); err != nil {
			log.Printf("Failed to save session: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error saving session"})
			return
		}
		type CustomUserResponse struct {
			UserID    uint      `json:"userId"`
			FirstName string    `json:"firstName"`
			LastName  string    `json:"lastName"`
			Email     string    `json:"email"`
			AuthType  string    `json:"authType"`
			Latitude  float64   `json:"latitude"`
			Longitude float64   `json:"longitude"`
			Address   *string   `json:"address"`
			CreatedAt time.Time `json:"createdAt"`
		}
		// Custom response; modify as needed. 
		response := CustomUserResponse{
			UserID:    foundUser.UserID,
			FirstName: foundUser.Entity.FirstName,
			LastName:  foundUser.Entity.LastName,
			Email:     foundUser.Email,
			AuthType:  foundUser.AuthType,
			Latitude:  foundUser.Latitude,
			Address:   foundUser.Address,
			Longitude: foundUser.Longitude,
			CreatedAt: foundUser.Entity.CreatedAt,
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"user":  response,
			"token": token,
		})
	}
}
