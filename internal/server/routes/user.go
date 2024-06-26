// This file contains the routes for the user model
//
// The routes here are as follows:
// - UserRoutes
// - .. more to be added later

package routes

import (
	"waitress-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRoutes sets up the routes for the user endpoints
func UserRoutes(router *gin.Engine, db *gorm.DB) {
	user := router.Group("api/users")
	{
		user.POST("/create", handlers.CreateUser(db))
		user.POST("/get", handlers.GetUser(db))
		user.POST("/update-user-location", handlers.UpdateUserLocation(db))
		user.POST("/update-account-info", handlers.UpdateUserAccountInformation(db))
		// user.POST("/", handlers.CreateUser)
		// user.GET("/:id", handlers.GetUser)
		// user.PUT("/:id", handlers.UpdateUser)
		// user.DELETE("/:id", handlers.DeleteUser)
	}
}
