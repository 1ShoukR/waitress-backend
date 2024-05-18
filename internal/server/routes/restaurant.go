package routes

import (
	"waitress-backend/internal/handlers"
	"waitress-backend/internal/utilities"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RestaurantRoutes(router *gin.Engine, db *gorm.DB) {
	userGroups := utilities.NewUserGroups()           // Initialize your user groups
	authGroups := utilities.NewAuthGroups(userGroups) // Create the auth groups from user groups

	restaurantRoutes := router.Group("api/restaurant")
	{
		// Use AuthGroups to apply middleware
		restaurantRoutes.POST("/create", utilities.UserRequired(authGroups, "Admin", "all"), handlers.CreateRestaurant(db, router))
		restaurantRoutes.POST("/edit", utilities.UserRequired(authGroups, "Admin", "all"), handlers.EditRestaurant(db, router))

		// Applying more specific or different groups as needed
		restaurantRoutes.POST("/local", utilities.UserRequired(authGroups, "Customer", "all"), handlers.GetLocalRestaurants(db, router))
		restaurantRoutes.POST("/reservations/:restaurantId/get", utilities.UserRequired(authGroups, "Staff", "all"), handlers.GetReservations(db, router))
		restaurantRoutes.GET("/avgrating/:restaurantId", utilities.UserRequired(authGroups, "Customer", "all"), handlers.GetAvgRating(db, router))
		restaurantRoutes.GET("/top10restaurants/", utilities.UserRequired(authGroups, "Customer", "all"), handlers.GetGlobalTopRestaurants(db, router))
	}
}
