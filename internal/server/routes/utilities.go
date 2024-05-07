package routes

import (
    "github.com/gin-gonic/gin"
    "waitress-backend/internal/handlers"
    "gorm.io/gorm"
)


func UtilitiesRoutes(router *gin.Engine, db *gorm.DB) {
	database := router.Group("api/db")
	{
		database.GET("/seed", handlers.Seed(db))
		database.GET("/run-all", handlers.RunAll(db))
	}
}