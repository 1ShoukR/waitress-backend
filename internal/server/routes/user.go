// This file contains the routes for the user model
//
// The routes here are as follows:
// - UserRoutes
// - .. more to be added later

package routes

import (
	"waitress-backend/internal/handlers"
	"waitress-backend/internal/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRoutes sets up the routes for the user endpoints
func UserRoutes(router *gin.Engine, db *gorm.DB) {
	user := router.Group("api/users")
	userGroups := utilities.NewUserGroups()           // Initialize your user groups
	authGroups := utilities.NewAuthGroups(userGroups) // Create the auth groups from user groups
	{
		user.POST("/create", handlers.CreateUser(db))
		user.POST("/get", handlers.GetUser(db))
		user.POST("/update-user-location", handlers.UpdateUserLocation(db))
		user.POST("/update-account-info", handlers.UpdateUserAccountInformation(db))
		user.POST("/get-admin-data", utilities.UserRequired(authGroups, "Admin", "all"), handlers.GetAdminData(db))
	}
}
