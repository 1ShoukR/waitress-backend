package routes

import (
	"github.com/gin-gonic/gin"
	"waitress-backend/internal/handlers"
	"gorm.io/gorm"
)


func RestaurantRoutes(router *gin.Engine, db *gorm.DB) {
	restaurantRoutes := router.Group("api/restaurant")
	{
		restaurantRoutes.POST("/create", handlers.CreateRestaurant(db, router))
		restaurantRoutes.POST("/edit", handlers.EditRestaurant(db, router))
		restaurantRoutes.POST("/local", handlers.GetLocalRestaurants(db, router))
		restaurantRoutes.POST("/reservations/:restaurantId/get", handlers.GetReservations(db, router))
	}
} 