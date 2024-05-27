package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	// "waitress-backend/internal/handlers"
	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid form data"})
			return
		}
		for name, value := range c.Request.PostForm {
			fmt.Println("name:", name)
			fmt.Println(name, value)
		}
		email := c.PostForm("email")
		password := c.PostForm("password")
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		userType := c.PostForm("userType")
		latitude := c.PostForm("latitude")
		longitude := c.PostForm("longitude")
		address := c.PostForm("address")
		city := c.PostForm("city")
		state := c.PostForm("state")
		zip := c.PostForm("zip")
		if email == "" || password == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid login credentials"})
			return
		}
		fmt.Println(email, password, firstName, lastName, userType, latitude, longitude, address, city, state, zip)
		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
			return
		}
		// Concatenate the address fields to create a single address string
		var newUser models.User
		fullUserAddress := address + ", " + city + ", " + state + " " + zip
		newUser.Email = email
		newUser.PasswordHash = string(hashedPassword)
		newUser.Entity.FirstName = firstName
		newUser.Entity.LastName = lastName
		newUser.AuthType = userType
		newUser.Latitude, _ = strconv.ParseFloat(latitude, 64)
		newUser.Longitude, _ = strconv.ParseFloat(longitude, 64)
		newUser.Address = &fullUserAddress

		// Create the user
		if err := db.Create(&newUser).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
			return
		}
		token, err := createToken(newUser.Email)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating token"})
			return
		}
		if err := verifyToken(token); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error verifying token"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"user":    newUser,
			"token":   token,
	})
	}
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