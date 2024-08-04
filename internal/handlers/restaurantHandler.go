// This file contains the handlers for the restaurant endpoints
//
// The handlers here are as follows:
// - EditRestaurant
// - GetLocalRestaurants
// - CreateRestaurant
// - GetReservations
// - GetSingleRestaurant
// - GetAvgRating
// - GetGlobalTopRestaurants

package handlers

import (
	// "log"
	// "errors"
	// "bytes"
	"fmt"
	// "io"
	"net/http"
	"strconv"
	"waitress-backend/internal/models"

	// "waitress-backend/internal/utilities"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "waitress-backend/internal/utilities"
)

// EditRestaurant is a handler for editing a restaurant
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
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ApiToken  string  `json:"apiToken"`
}

// GetLocalRestaurants is a handler for getting local restaurants based on user location
func GetLocalRestaurants(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurants []models.Restaurant
		var locationReq LocationRequest

		// Bind form data
		if err := c.ShouldBind(&locationReq); err != nil {
			fmt.Println("Error binding form data:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		userLat := locationReq.Latitude
		userLong := locationReq.Longitude
		apiToken := locationReq.ApiToken

		fmt.Println("Received apiToken:", apiToken)

		maxDistance := 50000.0 // Max distance in meters, increase for testing. Will need to be dynamic based on user input.

		// SQL query to calculate distance and filter restaurants
		query := `
			SELECT *, (
				6371000 * acos(
					cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
					sin(radians(?)) * sin(radians(latitude))
				)
			) AS distance
			FROM restaurant
			HAVING distance < ?
			ORDER BY distance
		`
		// Use raw SQL query to get nearby restaurants
		err := db.Raw(query, userLat, userLong, userLat, maxDistance).Scan(&restaurants).Error
		if err != nil {
			fmt.Println("Error executing the query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching nearby restaurants"})
			return
		}

		// Get the IDs of the filtered restaurants
		var restaurantIDs []uint
		for _, restaurant := range restaurants {
			restaurantIDs = append(restaurantIDs, restaurant.RestaurantId)
		}

		// Retrieve restaurants with preloaded ratings based on the filtered restaurant IDs
		err = db.Preload("Ratings").Preload("Categories").Where("restaurant_id IN (?)", restaurantIDs).Find(&restaurants).Error
		if err != nil {
			fmt.Println("Error preloading ratings:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload ratings"})
			return
		}

		if len(restaurants) == 0 {
			fmt.Println("No nearby restaurants found within", maxDistance, "meters.")
		} else {
			fmt.Println("Nearby restaurants found:", len(restaurants))
		}
		c.IndentedJSON(http.StatusOK, restaurants)
	}
}

// CreateRestaurant is a handler for creating a new restaurant in the database
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
			Name:           name,
			Address:        restaurantAddress,
			Phone:          restaurantPhone,
			Email:          restaurantEmail,
			Website:        &restaurantWebsite,
			NumberOfTables: &numberOfTables,
			Latitude:       &lat,
			Longitude:      &long,
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

// GetReservations is a handler for getting reservations for a specific restaurant
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

// GetSingleRestaurant is a handler for getting a single restaurant by ID
func GetSingleRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantId := c.Param("restaurantId")
		fmt.Printf("Restaurant ID: %s\n", restaurantId)
		var restaurant models.Restaurant
		results := db.Preload("Ratings").Preload("Categories").Preload("MenuItems").First(&restaurant, restaurantId)
		if results.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": results.Error.Error()})
			return
		}
		fmt.Printf("results: %+v\n", restaurant)
		c.IndentedJSON(http.StatusOK, restaurant)
	}
}

// Get's an individual menu item by ID
func GetMenuItem(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		menuItemId := c.Param("menuItemId")
		fmt.Printf("MenuItem ID: %s\n", menuItemId)
		var menuItem models.MenuItem
		results, err := menuItem.GetMenuItem(db, menuItemId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
			return
		}
		fmt.Printf("results: %+v\n", results)
		c.IndentedJSON(http.StatusOK, gin.H{"MenuItem": results})
	}
}

// GetAvgRating is a handler for getting the average rating of a restaurant
func GetAvgRating(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var avgRating float32
		id := c.Param("restaurantId")

		// Ideally this handler should just be responsible for only retrieving data. When ratings is more fleshed out r.CalcAvgRating and r.UpdateRating
		// should be called in endpoints where data is altered
		// -------------------------
		var r models.Restaurant
		avgRating, err := r.CalcAvgRating(db, id)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Problem calculating avg rating", "RestaurantID": id})
			return
		}

		err = r.UpdateAvgRating(db, id, avgRating)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Problem updating avg rating", "RestaurantID": id})
			return
		}
		// ---Remove this later--

		err = db.Table("restaurant").Select("average_rating").Where("restaurant_id = ?", id).Row().Scan(&avgRating)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Message": "Problem getting avg rating", "RestaurantID": id})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"RestaurantId": id, "AverageRating": avgRating})
	}
}

// GetGlobalTopRestaurants is a handler for getting the top 10 global restaurants based on average rating
func GetGlobalTopRestaurants(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurants []models.Restaurant
		err := db.Table("restaurant"). // Correct table name should be "restaurants"
			Preload("Ratings").
			Order("average_rating DESC").
			Limit(10).
			Preload("Categories").
			Preload("MenuItems"). // Ensure that your model has a MenuItems relation defined
			Find(&restaurants).Error
		if err != nil {
			fmt.Println("Error executing the query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching top restaurants"})
			return
		}

		c.IndentedJSON(http.StatusOK, restaurants)
	}
}
func GetAllCategories(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		result := db.Find(&categories)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories", "message": result.Error.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, categories)
	}
}