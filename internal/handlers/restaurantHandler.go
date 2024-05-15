package handlers

import (
	// "log"
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"waitress-backend/internal/models"
	"waitress-backend/internal/utilities"

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

func GetLocalRestaurants(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO - Make this based off of the user's location
		var restaurants []models.Restaurant
		userLatStr := c.PostForm("latitude")
		userLongStr := c.PostForm("longitude")
		apiToken := c.PostForm("apiToken")

		fmt.Println("Received latitude:", userLatStr)
		fmt.Println("Received longitude:", userLongStr)
		fmt.Println("Received apiToken:", apiToken)

		userLat, err := strconv.ParseFloat(userLatStr, 64)
		if err != nil {
			fmt.Println("Error parsing latitude:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}

		userLong, err := strconv.ParseFloat(userLongStr, 64)
		if err != nil {
			fmt.Println("Error parsing longitude:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		maxDistance := 5000.0 // Max distance in meters

		// Retrieve all restaurants
		// TODO: We need to figure a way to get restaurants based on the user's location rather than all restaurants
		if err := db.Find(&restaurants).Error; err != nil {
			fmt.Println("Error retrieving restaurants:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
			return
		}

		if len(restaurants) == 0 {
			fmt.Println("No restaurants found in database.")
		} else {
			fmt.Println("Total restaurants retrieved:", len(restaurants))
		}

		var nearbyRestaurants []models.Restaurant
		// Need to figure out how to get the distance between the user and the restaurant
		// and only return the restaurants that are within the max distance
		for _, restaurant := range restaurants {
			if utilities.Haversine(userLat, userLong, *restaurant.Latitude, *restaurant.Longitude) <= maxDistance {
				nearbyRestaurants = append(nearbyRestaurants, restaurant)
			}
		}

		if len(nearbyRestaurants) == 0 {
			fmt.Println("No nearby restaurants found within", maxDistance, "meters.")
		} else {
			fmt.Println("Nearby restaurants found:", len(nearbyRestaurants))
		}

		c.JSON(http.StatusOK, gin.H{"restaurants": nearbyRestaurants})
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
		reservationId := c.Param("restaurantId")
		fmt.Printf("Reservation ID: %s\n", reservationId)
		var reservationList []models.Reservation
		results := db.Find(&reservationList)
		if results.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": results.Error.Error()})
		}
		c.IndentedJSON(http.StatusOK, reservationList)
	}
}