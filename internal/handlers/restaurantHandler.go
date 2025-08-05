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
// - GetAvailableTables

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
	"waitress-backend/internal/utilities"

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
		fmt.Printf("User location: Latitude: %f, Longitude: %f\n", userLat, userLong)

		maxDistance := 100000.0 // Max distance in meters, increase for testing

		// SQL query to calculate distance and filter restaurants
		query := `
			SELECT *, (
				6371000 * acos(
					cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
					sin(radians(?)) * sin(radians(latitude))
				)
			) AS distance
			FROM restaurants
			HAVING distance < ?
			ORDER BY distance
		`
		fmt.Printf("Query parameters: userLat=%f, userLong=%f, maxDistance=%f\n", userLat, userLong, maxDistance)
		// Use raw SQL query to get nearby restaurants
		err := db.Raw(query, userLat, userLong, userLat, maxDistance).Scan(&restaurants).Error
		if err != nil {
			fmt.Println("Error executing the query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching nearby restaurants"})
			return
		}

		// Debugging: Output the result of the distance calculation
		for _, restaurant := range restaurants {
			fmt.Printf("Restaurant: %s, Latitude: %f, Longitude: %f,\n",
				restaurant.Name, *restaurant.Latitude, *restaurant.Longitude)
		}

		// Return nearby restaurants
		c.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
	}
}

// CreateRestaurant is a handler for creating a restaurant
func CreateRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurant models.Restaurant

		// Bind JSON to Restaurant struct
		if err := c.ShouldBindJSON(&restaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "message": err.Error()})
			return
		}

		// Insert the restaurant into the database
		if err := db.Create(&restaurant).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating restaurant", "message": err.Error()})
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
	}
}

type ReservationRequest struct {
	ApiToken string `json:"apiToken"`
}

// GetReservations is a handler for getting reservations for a restaurant
func GetReservations(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reservations []models.Reservation
		var resRequest ReservationRequest

		// Bind JSON data
		if err := c.ShouldBind(&resRequest); err != nil {
			fmt.Println("Error binding form data:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		apiToken := resRequest.ApiToken

		fmt.Println("Received apiToken:", apiToken)

		// Query for reservations - you'll need to implement the filtering logic here
		err := db.Find(&reservations).Error
		if err != nil {
			fmt.Println("Error executing the query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations"})
			return
		}

		// Return reservations
		c.JSON(http.StatusOK, gin.H{"reservations": reservations})
	}
}

// GetSingleRestaurant is a handler for getting a single restaurant based on ID
func GetSingleRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the restaurant ID from the URL parameter
		restaurantIDStr := c.Param("restaurantId") // Fixed: match route parameter name
		restaurantID, err := strconv.Atoi(restaurantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		var restaurant models.Restaurant

		// Query for the restaurant by ID with all necessary associations
		err = db.Preload("MenuItems").Preload("Categories").Preload("Ratings").First(&restaurant, restaurantID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurant"})
			}
			return
		}

		// Return the restaurant
		c.JSON(http.StatusOK, gin.H{"restaurant": restaurant})
	}
}

func GetAvgRating(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract restaurant ID from the URL parameter
		restaurantIDStr := c.Param("id")

		// Get the average rating for the restaurant with the given ID
		var restaurant models.Restaurant
		rating, err := restaurant.CalcAvgRating(db, restaurantIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average rating", "message": err.Error()})
			return
		}

		// Return the average rating
		c.JSON(http.StatusOK, gin.H{"average_rating": rating})
	}
}

// GetGlobalTopRestaurants is a handler for getting top-rated restaurants globally
func GetGlobalTopRestaurants(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var restaurants []models.Restaurant

		// Query for top-rated restaurants - adjust the query as needed
		// For example, you might want to order by rating, review count, etc.
		err := db.Limit(10).Find(&restaurants).Error // Limit to top 10 for performance
		if err != nil {
			fmt.Println("Error executing the query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching top restaurants"})
			return
		}

		// Return top restaurants
		c.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
	}
}

// UserToFavorites updates a user's favorites list to add/remove a restaurant
func UserToFavorites(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get restaurant ID from the URL parameter
		restaurantIDStr := c.Param("id")
		restId, err := strconv.Atoi(restaurantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}

		// Get user ID from the POST form data
		userIDStr := c.PostForm("userId")
		userId, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Get the user
		var user models.User
		err = db.First(&user, userId).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		err = user.AddToFavorites(db, uint(restId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to favorites", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant added to favorites"})
	}
}

// GetAvailableTables handles the enhanced table selection API with comprehensive filtering
func GetAvailableTables(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract restaurant ID from URL parameter - following established pattern
		restaurantIDStr := c.Param("restaurantId")
		restaurantID, err := strconv.ParseUint(restaurantIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID format"})
			return
		}

		// Validate restaurant exists
		var restaurant models.Restaurant
		if err := db.First(&restaurant, uint(restaurantID)).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating restaurant"})
			}
			return
		}

		// Validate and parse query parameters using our utility
		filters, err := utilities.ValidateTableFilters(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build dynamic query for available tables
		query := utilities.BuildAvailableTableQuery(db, uint(restaurantID), filters)

		// Execute query to get matching tables
		var tables []models.Table
		if err := query.Find(&tables).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tables"})
			return
		}

		// Convert to API response format optimized for React Native
		var results []utilities.TableQueryResult
		for _, table := range tables {
			results = append(results, utilities.ConvertTableToResult(table))
		}

		// Build comprehensive response
		response := gin.H{
			"restaurant": gin.H{
				"id":   restaurant.RestaurantId,
				"name": restaurant.Name,
			},
			"filters": gin.H{
				"zone":        filters.Zone,
				"tableType":   filters.TableType,
				"view":        filters.View,
				"minCapacity": filters.MinCapacity,
				"maxCapacity": filters.MaxCapacity,
				"available":   filters.Available,
			},
			"tables": results,
			"count":  len(results),
		}

		c.JSON(http.StatusOK, response)
	}
}
