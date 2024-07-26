// This file contains the handlers for user related endpoints 
//
// The handlers here are as follows:
// - CreateUser
// - GetUser
// - UpdateUserLocation
// - UpdateUserAccountInformation


package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	// "waitress-backend/internal/handlers"
	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser is a handler for creating a new user account
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	UserType  string  `json:"userType"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data", "error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
			return
		}

		var newUser models.User
		fullUserAddress := req.Address + ", " + req.City + ", " + req.State + " " + req.Zip
		newUser.Email = req.Email
		newUser.PasswordHash = string(hashedPassword)
		newUser.Entity.FirstName = req.FirstName
		newUser.Entity.LastName = req.LastName
		newUser.AuthType = req.UserType
		newUser.Latitude = req.Latitude
		newUser.Longitude = req.Longitude
		newUser.Address = &fullUserAddress

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
			return
		}

		token, err := createToken(newUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating token"})
			return
		}
		if err := verifyToken(token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error verifying token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"user":    newUser,
			"token":   token,
		})
	}
}
// GetUser is a handler for getting all users in the database
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

// UpdateUserLocation is a handler for updating a user's location based on their user latitude and longitude
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

// UpdateUserAccountInformation is a handler for updating a user's account information from the Edit Account page
func UpdateUserAccountInformation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.User
		userId := c.PostForm("userId")
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		street := c.PostForm("street")
		city := c.PostForm("city")
		state := c.PostForm("state")
		zip := c.PostForm("zip")
		address := street + ", " + city + ", " + state + " " + zip

		result := db.Preload("Entity").Where("user_id = ?", userId).First(&foundUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			}
			return
		}

		// Update user information
		updatedUser, err := foundUser.UpdateAccountInformation(db, firstName, lastName, email, address, city, state, zip, phone); 
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		token := sessions.Default(c).Get("apiToken")
		c.JSON(http.StatusOK, gin.H{"message": "User account information updated successfully", "user": updatedUser, "token": token})
	}
}