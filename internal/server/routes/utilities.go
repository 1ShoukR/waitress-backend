package routes

import (
	"waitress-backend/internal/handlers"
	"waitress-backend/internal/utilities"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UtilitiesRoutes(router *gin.Engine, db *gorm.DB) {
	userGroups := utilities.NewUserGroups()           // Initialize your user groups
	authGroups := utilities.NewAuthGroups(userGroups) // Create the auth groups from user groups
	database := router.Group("api/db")
	{
		database.GET("/seed", utilities.UserRequired(authGroups, "Dev", "all"), handlers.Seed(db))
		database.GET("/run-all", handlers.RunAll(db))
		database.GET("/migrate", handlers.MigrateDb(db))
	}
}
