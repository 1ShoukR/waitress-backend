package routes

import (
    "github.com/gin-gonic/gin"
    "waitress-backend/internal/handlers"
    "gorm.io/gorm"
	"waitress-backend/internal/utilities"
)


func UtilitiesRoutes(router *gin.Engine, db *gorm.DB) {
	userGroups := utilities.NewUserGroups() // Initialize your user groups
	authGroups := utilities.NewAuthGroups(userGroups) // Create the auth groups from user groups
	database := router.Group("api/db")
	{
		database.GET("/seed", utilities.UserRequired(authGroups, "Dev", "all"), handlers.Seed(db))
		database.GET("/run-all", utilities.UserRequired(authGroups, "Dev", "all"), handlers.RunAll(db))
	}
}