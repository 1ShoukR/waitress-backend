package handlers


import (
	// "log"
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"waitress-backend/internal/models"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "waitress-backend/internal/utilities"
)

func EditRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here, we will query for a restaurant based 
		// On the owner/admin. Using a discriminator,
		// We can control what the flow of this func
		// looks like. For example, we can edit the tables,
		// menu, employees, available tables? etc.
		// delimiter := c.PostForm("delimiter")
		// fmt.printf("%+v\n",delimiter)

	}
}

func CreateRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse and validate form inputs
		name := c.PostForm("restaurantName")
		restaurantAddress := c.PostForm("restaurantAddress")
		restaurantPhone := c.PostForm("restaurantPhone")
		restaurantEmail := c.PostForm("restaurantEmail")
		restaurantWebsite := c.PostForm("restaurantWebsite")
		restaurantNumberOfTablesStr := c.PostForm("restaurantNumberOfTables")
		restaurantLatStr := c.PostForm("restaurantLat")
		restaurantLongStr := c.PostForm("restaurantLong")
		// Convert restaurantNumberOfTables to int
		numberOfTables, err := strconv.Atoi(restaurantNumberOfTablesStr)
		if err != nil {
			// Handle error, perhaps set an HTTP error response
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid number of tables"})
			return
		}
		// Convert restaurantLat to float64
		lat, err := strconv.ParseFloat(restaurantLatStr, 64)
		if err != nil {
			// Handle error
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
			return
		}
		// Convert restaurantLong to float64
		long, err := strconv.ParseFloat(restaurantLongStr, 64)
		if err != nil {
			// Handle error
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
			return
		}
		// Create newRestaurant struct with correct types
		newRestaurant := models.Restaurant{
			Name: name,
			Address: restaurantAddress,
			Phone: restaurantPhone,
			Email: restaurantEmail,
			Website: &restaurantWebsite,
			NumberOfTables: &numberOfTables,
			Latitude: &lat,
			Longitude: &long,
		}
		// Print the struct for debugging
		fmt.Printf("%+v\n", newRestaurant)
		// Save to database
		result := db.Create(&newRestaurant)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create restaurant"})
			return
		}
		// Respond with success or forward the error
		c.JSON(http.StatusOK, gin.H{"message": "restaurant created"})
	}
}


func GetReservations(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// We are going to transform this to be a reservation from 
		// a restaurant, based on a user. 
		var reservationList []models.Reservation
		results := db.Find(&reservationList)
		if results.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": results.Error.Error()})
		}
		c.IndentedJSON(http.StatusOK, reservationList)
	}
}